package admin

import (
	"clanplatform/internal/db"
	"clanplatform/internal/services"
)

type storage interface {
	ListUsers() ([]db.User, error)
	CreateUser(user db.User) (*db.User, error)

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
	CreateProductVariant(productID int64, title string, inventory int) (*db.ProductVariant, error)
	GetProductByID(id int64) (*db.ProductWithDetails, error)
	CreateProduct(product db.Product) (*db.Product, error)
	UpdateProduct(id int64, update map[string]any) (*db.Product, error)
}

type emailClient interface {
	SendEmail(message *services.MailMessage) error
}

type admin struct {
	storage     storage
	emailClient emailClient
}

func NewAdmin(storage storage, emailClient emailClient) *admin {
	return &admin{
		storage:     storage,
		emailClient: emailClient,
	}
}
