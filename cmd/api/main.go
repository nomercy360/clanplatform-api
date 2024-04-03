package main

import (
	"clanplatform/internal/api"
	"clanplatform/internal/db"
	"clanplatform/internal/services"
	"flag"
	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

// Config represents the configuration structure
type Config struct {
	DbConnString string       `yaml:"db_conn_string"`
	Server       ServerConfig `yaml:"server"`
	JWTSecret    string       `yaml:"jwt_secret"`
	ResendApiKey string       `yaml:"resend_api_key"`
}

// ServerConfig represents the server configuration structure
type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
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

	return &config, nil
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the config file")
	flag.Parse()

	config, err := ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	storage, err := db.NewDB(config.DbConnString)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}

	defer storage.Close()

	// Create a new Echo instance
	r := chi.NewRouter()

	// Create a new API instance

	email := services.NewEmailClient(config.ResendApiKey)

	app := api.New(storage, email, config.JWTSecret)

	app.RegisterRoutes(r)

	// Start the server
	log.Printf("Starting server on %s:%s\n", config.Server.Host, config.Server.Port)

	if err := http.ListenAndServe(config.Server.Host+":"+config.Server.Port, r); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
