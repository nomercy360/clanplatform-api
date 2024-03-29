package api

import "clanplatform/internal/entity"

type storage interface {
	ListUsers() ([]entity.User, error)
	InviteUser(token string, email string, enum entity.UserRoleEnum) error
	CreateUser(user entity.User) (entity.User, error)
	ListInvites() ([]entity.Invite, error)
	GetInviteByEmail(email string) (*entity.Invite, error)
	GetUserByEmail(email string) (entity.User, error)

	Ping() error

	CreateDiscount(discount entity.Discount) (entity.Discount, error)
	GetDiscounts() ([]entity.Discount, error)
	UpdateDiscount(discount entity.Discount) (entity.Discount, error)
	DeleteDiscount(id string) error
}
