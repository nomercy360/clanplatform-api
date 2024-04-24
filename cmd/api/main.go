package main

import (
	adm "clanplatform/internal/admin"
	"clanplatform/internal/db"
	"clanplatform/internal/services"
	storeSvc "clanplatform/internal/store"
	"clanplatform/internal/terrors"
	"clanplatform/internal/transport"
	"context"
	"errors"
	"flag"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/yaml.v3"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	DbConnString string       `yaml:"db_conn_string" validate:"required"`
	Server       ServerConfig `yaml:"server" validate:"required"`
	JWTSecret    string       `yaml:"jwt_secret" validate:"required"`
	ResendApiKey string       `yaml:"resend_api_key" validate:"required"`
}

type ServerConfig struct {
	Port string `yaml:"port" validate:"required"`
	Host string `yaml:"host" validate:"required"`
}

// ReadConfig reads and unmarshal the YAML configuration from the given file
func ReadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the config file")
	flag.Parse()

	config, err := ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	pg, err := db.NewDB(config.DbConnString)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}

	defer pg.Close()

	e := echo.New()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var (
			code = http.StatusInternalServerError
			msg  interface{}
		)

		var he *echo.HTTPError
		var terror *terrors.Error
		if errors.As(err, &he) {
			code = he.Code
			msg = he.Message
		} else if errors.As(err, &terror) {
			code = terror.Code
			msg = terror.Message
		} else {
			msg = err.Error()
		}

		if _, ok := msg.(string); ok {
			msg = map[string]interface{}{"error": msg}
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, msg)
			}

			if err != nil {
				e.Logger.Error(err)
			}
		}
	}

	email := services.NewEmailClient(config.ResendApiKey)

	admin := adm.NewAdmin(pg, email)

	store := storeSvc.NewStore(pg)

	app := transport.NewTransport(admin, store, config.JWTSecret)

	app.RegisterRoutes(e)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
