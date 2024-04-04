package db

import "time"

type User struct {
	ID           int64        `db:"id" json:"id"`
	Email        string       `db:"email" json:"email"`
	PasswordHash string       `db:"password_hash" json:"-"`
	FirstName    string       `db:"first_name" json:"first_name"`
	LastName     string       `db:"last_name" json:"last_name"`
	Role         UserRoleEnum `db:"role" json:"role"`
	CreatedAt    time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time   `db:"deleted_at" json:"-"`
}

func (s *storage) ListUsers() ([]User, error) {
	users := make([]User, 0)

	err := s.pg.Select(&users, "SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storage) CreateUser(user User) (*User, error) {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, role)
		VALUES (:email, :password_hash, :first_name, :last_name, :role)
		RETURNING id, email, first_name, last_name, role, created_at, updated_at, deleted_at;
	`

	rows, err := s.pg.NamedQuery(query, user)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *storage) GetUserByEmail(email string) (*User, error) {
	var user User

	query := `
		SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1;
	`

	err := s.pg.Get(&user, query, email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
