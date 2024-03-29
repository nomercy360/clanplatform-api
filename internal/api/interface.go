package api

import "clanplatform/internal/entity"

type storage interface {
	ListUsers() ([]entity.User, error)
	Ping() error
}
