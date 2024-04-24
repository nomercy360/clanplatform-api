package admin

import (
	"clanplatform/internal/db"
	"clanplatform/internal/terrors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (adm *admin) ListUsers() ([]db.User, error) {
	users, err := adm.storage.ListUsers()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
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
} // @Name UserWithToken

type AuthUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
} // @Name AuthUser

func (adm *admin) AuthUser(authUser AuthUser) (*UserWithToken, error) {
	user, err := adm.storage.GetUserByEmail(authUser.Email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(authUser.Password)); err != nil {
		return nil, terrors.BadRequest(err)
	}

	claims := jwt.MapClaims{
		"email": user.Email,
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
} // @Name CreateUser

func (adm *admin) CreateUser(cr CreateUser) (*db.User, error) {
	hashedPassword, _ := hashPassword(cr.Password)

	user := db.User{
		Email:        cr.Email,
		PasswordHash: hashedPassword,
		FullName:     cr.FullName,
	}

	createdUser, err := adm.storage.CreateUser(user)

	if err != nil && db.IsDuplicationError(err) {
		return nil, terrors.BadRequest(err)
	} else if err != nil {
		return nil, err
	}

	return createdUser, nil
}
