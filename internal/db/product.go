package db

import (
	"clanplatform/internal/entity"
)

func (s *Storage) ListProducts() ([]entity.User, error) {
	users := make([]entity.User, 0)

	err := s.pg.Select(&users, "SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	return users, nil
}
