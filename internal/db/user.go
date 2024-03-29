package db

import (
	"clanplatform/internal/entity"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	pg *sqlx.DB
}

func (s *Storage) ListUsers() ([]entity.User, error) {
	users := make([]entity.User, 0)

	err := s.pg.Select(&users, "SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	return users, nil
}
