package transport

import (
	adm "clanplatform/internal/admin"
	"clanplatform/internal/db"
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
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
	store     store
	jwtSecret string
}

type admin interface {
	ListUsers() ([]db.User, error)
	CreateUser(user adm.CreateUser) (*db.User, error)
	AuthUser(email, password string) (*adm.UserWithToken, error)
	GetUserByEmail(email string) (*db.User, error)

	CreateDiscount(cd adm.CreateDiscount) (*db.Discount, error)
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
	CreateProduct(product db.Product) (*db.Product, error)
	UpdateProduct(id int64, update adm.UpdateProductRequest) (*db.Product, error)
}

type store interface {
	ListProducts() ([]db.ProductWithDetails, error)
	GetProductByID(id int64) (*db.ProductWithDetails, error)
}

func NewTransport(admin admin, store store, jwtSecret string) *transport {
	return &transport{admin: admin, store: store, jwtSecret: jwtSecret}
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

func (tr *transport) HealthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthStatus{Status: "ok"})
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (tr *transport) RegisterRoutes(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/health", tr.HealthCheckHandler)

	a := e.Group("/admin")
	a.POST("/auth", tr.AuthCookieHandler)
	a.POST("/auth/token", tr.AuthTokenHandler)
	a.POST("/auth", tr.AuthCookieHandler)
	a.GET("/users", tr.ListUsersHandler)
	a.POST("/users", tr.CreateUserHandler)

	a.GET("/discounts", tr.ListDiscountsHandler)
	a.POST("/discounts", tr.CreateDiscountHandler)

	a.GET("/collections", tr.ListCollectionsHandler)
	a.POST("/collections", tr.CreateCollectionHandler)
	a.GET("/collections/{id}", tr.GetCollectionByIDHandler)
	a.POST("/collections/{id}", tr.UpdateCollectionHandler)
	a.DELETE("/collections/{id}", tr.DeleteCollectionHandler)
	a.POST("/collections/{id}/products", tr.AddProductsToCollectionHandler)
	a.DELETE("/collections/{id}/products", tr.RemoveProductsFromCollectionHandler)

	a.GET("/products", tr.ListProductsHandler)
	a.POST("/products/:id/variants", tr.CreateProductVariantHandler)
	a.GET("/products/:id", tr.GetProductByIDHandler)
	a.POST("/products", tr.CreateProductHandler)
	a.POST("/products/{id}", tr.UpdateProductHandler)

	st := e.Group("/store")

	st.GET("/products", tr.ListProductsHandler)
	st.GET("/products/{id}", tr.GetProductByIDHandler)

}
