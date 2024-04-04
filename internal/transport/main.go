package transport

import (
	admSvc "clanplatform/internal/admin"
	"clanplatform/internal/db"
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

type transport struct {
	admin     admin
	jwtSecret string
}

type admin interface {
	ListUsers() ([]db.User, error)
	CreateUser(user db.User) (*db.User, error)
	AuthUser(email, password string) (*admSvc.UserWithToken, error)

	ListInvites() ([]db.Invite, error)
	GetInviteByEmail(email string) (*db.Invite, error)
	AcceptInvite(token, password, firstName, lastName string) error
	InviteUser(role, email string) error
	GetUserByEmail(email string) (*db.User, error)

	CreateDiscount(discount db.Discount) (*db.Discount, error)
	ListDiscounts() ([]db.Discount, error)
	UpdateDiscount(discount db.Discount) (*db.Discount, error)
	DeleteDiscount(id string) error

	ListCollections() ([]db.ProductCollection, error)
	CreateCollection(title string, handle string) (*db.ProductCollection, error)
	GetCollectionByID(id int64) (*db.ProductCollection, error)
	UpdateCollection(title *string, handle *string, id int64) (*db.ProductCollection, error)
	DeleteCollection(id int64) error
	AddProductsToCollection(collectionID int64, productIDs []int64) error
	RemoveProductsFromCollection(collectionID int64, productIDs []int64) error

	ListProducts() ([]db.ProductWithDetails, error)
	CreateProductVariant(variant db.ProductVariant) (*db.ProductVariant, error)
	GetProductByID(id int64) (*db.ProductWithDetails, error)
}

func New(admin admin, jwtSecret string) *transport {
	return &transport{admin: admin, jwtSecret: jwtSecret}
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

func (tr *transport) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	_ = WriteJSON(w, http.StatusOK, HealthStatus{Status: "ok"})
}

func (tr *transport) RegisterRoutes(r chi.Router) {
	r.Get("/health", tr.HealthCheckHandler)

	r.Mount("/admin", AdminRoutes(tr))
}

func AdminRoutes(tr *transport) http.Handler {
	r := chi.NewRouter()

	r.Use(WithAuth("secret"))

	r.Post("/auth", tr.AuthCookieHandler)
	r.Post("/auth/token", tr.AuthTokenHandler)
	r.Get("/users", tr.ListUsersHandler)

	r.Get("/discounts", tr.ListDiscountsHandler)
	r.Post("/discounts", tr.CreateDiscountHandler)

	r.Post("/invites", tr.InviteUserHandler)
	r.Post("/invites/accept", tr.AcceptInviteHandler)
	r.Get("/invites", tr.ListInvitesHandler)

	r.Get("/collections", tr.ListCollectionsHandler)
	r.Post("/collections", tr.CreateCollectionHandler)
	r.Get("/collections/{id}", tr.GetCollectionByIDHandler)
	r.Post("/collections/{id}", tr.UpdateCollectionHandler)
	r.Delete("/collections/{id}", tr.DeleteCollectionHandler)
	r.Post("/collections/{id}/products", tr.AddProductsToCollectionHandler)
	r.Delete("/collections/{id}/products", tr.RemoveProductsFromCollectionHandler)

	r.Get("/products", tr.ListProductsHandler)
	r.Post("/products/{id}/variants", tr.CreateProductVariantHandler)
	r.Get("/products/{id}", tr.GetProductByIDHandler)

	return r
}
