package api

import (
	"clanplatform/internal/entity"
	"clanplatform/internal/services"
)

type storage interface {
	ListUsers() ([]entity.User, error)
	CreateUser(user entity.User) (entity.User, error)

	ListInvites() ([]entity.Invite, error)
	GetInviteByEmail(email string) (*entity.Invite, error)
	InviteUser(token string, email string, enum entity.UserRoleEnum) error
	GetUserByEmail(email string) (entity.User, error)

	Ping() error

	CreateDiscount(discount entity.Discount) (entity.Discount, error)
	GetDiscounts() ([]entity.Discount, error)
	UpdateDiscount(discount entity.Discount) (entity.Discount, error)
	DeleteDiscount(id string) error

	ListCollections() ([]entity.ProductCollection, error)
	CreateCollection(title string, handle string) (entity.ProductCollection, error)
	GetCollectionByID(id int64) (entity.ProductCollection, error)
	UpdateCollection(title *string, handle *string, id int64) (entity.ProductCollection, error)
	DeleteCollection(id int64) error
	AddProductsToCollection(collectionID int64, productIDs []int64) error
	RemoveProductsFromCollection(collectionID int64, productIDs []int64) error

	ListProducts() ([]entity.Product, error)
	CreateProductVariant(productID int64, title string, priceList []entity.Price, inventory int) (entity.ProductVariant, error)
}

type emailService interface {
	SendEmail(message *services.MailMessage) error
}
