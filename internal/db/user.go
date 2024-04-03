package db

import (
	"clanplatform/internal/entity"
)

func (s *storage) ListUsers() ([]entity.User, error) {
	users := make([]entity.User, 0)

	err := s.pg.Select(&users, "SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storage) CreateUser(user entity.User) (entity.User, error) {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, role)
		VALUES (:email, :password_hash, :first_name, :last_name, :role);
	`

	_, err := s.pg.NamedExec(query, user)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *storage) GetUserByEmail(email string) (entity.User, error) {
	var user entity.User

	err := s.pg.Get(&user, "SELECT * FROM users WHERE email = $1", email)

	if err != nil {
		if IsNoRowsError(err) {
			return entity.User{}, nil
		} else {
			return entity.User{}, err
		}
	}

	return user, nil
}
