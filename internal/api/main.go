package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type api struct {
	storage storage
}

func New(storage storage) *api {
	return &api{storage: storage}
}

// WriteError responds to a HTTP request with an error.
func WriteError(w http.ResponseWriter, code int, message string) error {
	err := WriteJSON(w, code, map[string]string{"error": message})
	if err != nil {
		return err
	}

	return nil
}

// WriteJSON writes a JSON response to a HTTP request.
func WriteJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		return err
	}

	return nil
}

type HealthStatus struct {
	Status string `json:"status"`
}

func (api *api) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := api.storage.Ping()

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, "Database connection error")
		return
	}

	_ = WriteJSON(w, http.StatusOK, HealthStatus{Status: "ok"})
}

func (api *api) RegisterRoutes(r chi.Router) {
	r.Get("/health", api.HealthCheckHandler)

	r.Mount("/admin", AdminRoutes(api))
}

func AdminRoutes(api *api) http.Handler {
	r := chi.NewRouter()

	r.Use(WithAuth("secret"))

	r.Post("/auth", api.Auth)
	r.Post("/auth/token", api.Auth)
	r.Get("/users", api.ListUsers)
	r.Get("/discounts", api.ListDiscounts)
	r.Post("/discounts", api.CreateDiscount)
	r.Post("/invites", api.InviteUser)
	r.Post("/invites/accept", api.AcceptInvite)

	return r
}
