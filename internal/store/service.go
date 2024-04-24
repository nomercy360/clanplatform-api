package store

import (
	"clanplatform/internal/db"
	"clanplatform/internal/services"
)

type storage interface {
	ListProducts() ([]db.ProductWithDetails, error)
	GetProductByID(id int64) (*db.ProductWithDetails, error)
}

type emailClient interface {
	SendEmail(message *services.MailMessage) error
}

type store struct {
	storage storage
}

func NewStore(storage storage) *store {
	return &store{
		storage: storage,
	}
}
