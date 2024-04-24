package admin

import (
	"bytes"
	"clanplatform/internal/db"
	"clanplatform/internal/services"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"time"
)

func (adm *admin) ListUsers() ([]db.User, error) {
	users, err := adm.storage.ListUsers()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (adm *admin) InviteUser(role, email string) error {
	mySigningKey := []byte("secret")

	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		Subject:   email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(mySigningKey)

	err = adm.storage.InviteUser(ss, email)

	if err != nil {
		return err
	}

	type Context struct {
		InviteURL string
	}

	tmpl, err := template.ParseFiles("templates/otp.gohtml")

	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	// TODO: change to real URL when frontend is ready
	err = tmpl.Execute(&tpl, Context{InviteURL: "http://localhost:8080/invite/accept?token=" + ss})

	if err != nil {
		return err
	}

	message := &services.MailMessage{
		To:       email,
		Subject:  "You have been invited to join the platform",
		From:     "hi@mxksim.dev",
		HtmlBody: tpl.String(),
	}

	if err = adm.emailClient.SendEmail(message); err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (adm *admin) ListInvites() ([]db.Invite, error) {
	invites, err := adm.storage.ListInvites()

	if err != nil {
		return nil, err
	}

	return invites, nil
}

func (adm *admin) AcceptInvite(token, password, fullName string) error {
	if token == "" || password == "" {
		return invalidReqErr
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil || !parsed.Valid {
		return errors.New("invalid token")
	}

	email, _ := parsed.Claims.GetSubject()

	if email == "" {
		return errors.New("invalid token")
	}

	invite, _ := adm.storage.GetInviteByEmail(email)

	if invite == nil {
		return errors.New("invite not found")
	}

	hashedPassword, _ := hashPassword(password)

	user := db.User{
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if _, err = adm.storage.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (adm *admin) GetUserByEmail(email string) (*db.User, error) {
	user, err := adm.storage.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

type UserWithToken struct {
	User  *db.User `json:"user"`
	Token string   `json:"token"`
}

func (adm *admin) AuthUser(email, password string) (*UserWithToken, error) {
	user, err := adm.storage.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString([]byte("secret"))

	userWithToken := UserWithToken{
		User:  user,
		Token: tokenString,
	}

	return &userWithToken, nil
}

type CreateUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
}

func (adm *admin) CreateUser(cr CreateUser) (*db.User, error) {
	hashedPassword, _ := hashPassword(cr.Password)

	user := db.User{
		Email:        cr.Email,
		PasswordHash: hashedPassword,
		FullName:     cr.FullName,
	}

	createdUser, err := adm.storage.CreateUser(user)

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (adm *admin) GetInviteByEmail(email string) (*db.Invite, error) {
	invite, err := adm.storage.GetInviteByEmail(email)

	if err != nil {
		return nil, err
	}

	return invite, nil
}
