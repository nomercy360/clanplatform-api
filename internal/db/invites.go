package db

import (
	"time"
)

type Invite struct {
	ID        int64      `db:"id"`
	Email     string     `db:"email"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	ExpiresAt time.Time  `db:"expires_at"`
	Accepted  bool       `db:"accepted"`
	DeletedAt *time.Time `db:"deleted_at"`
	Token     string     `db:"token"`
}

func (s *Storage) InviteUser(token string, email string) error {
	invite := Invite{
		Token:     token,
		Email:     email,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	_, err := s.pg.NamedExec("INSERT INTO invites (email, role, expires_at, token) VALUES (:email, :role, :expires_at, :token)", invite)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListInvites() ([]Invite, error) {
	invites := make([]Invite, 0)

	err := s.pg.Select(&invites, "SELECT * FROM invites")

	if err != nil {
		return nil, err
	}

	return invites, nil
}

func (s *Storage) GetInviteByEmail(email string) (*Invite, error) {
	var invite Invite

	err := s.pg.Get(&invite, "SELECT * FROM invites WHERE email = $1", email)

	if err != nil {
		if IsNoRowsError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &invite, nil
}
