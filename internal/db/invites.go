package db

import (
	"clanplatform/internal/entity"
	"time"
)

func (s *storage) InviteUser(token string, email string, role entity.UserRoleEnum) error {
	invite := entity.Invite{
		Token:     token,
		Email:     email,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		Role:      role,
	}

	_, err := s.pg.NamedExec("INSERT INTO invites (email, role, expires_at, token) VALUES (:email, :role, :expires_at, :token)", invite)

	if err != nil {
		return err
	}

	return nil
}

func (s *storage) ListInvites() ([]entity.Invite, error) {
	var invites []entity.Invite

	err := s.pg.Select(&invites, "SELECT * FROM invites")

	if err != nil {
		return nil, err
	}

	return invites, nil
}

func (s *storage) GetInviteByEmail(email string) (*entity.Invite, error) {
	var invite entity.Invite

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
