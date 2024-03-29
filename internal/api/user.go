package api

import (
	"clanplatform/internal/entity"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Function to hash password.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (api *api) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := api.storage.ListUsers()

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, users)
}

func decodeCredentials(r *http.Request) (email, password string, err error) {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err = json.NewDecoder(r.Body).Decode(&data)
	return data.Email, data.Password, err
}

func (api *api) verifyUserCredentials(email, password string) (entity.User, error) {
	user, err := api.storage.GetUserByEmail(email)
	if err != nil {
		return entity.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return user, err
}

func (api *api) Auth(w http.ResponseWriter, r *http.Request) {
	email, password, err := decodeCredentials(r)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.verifyUserCredentials(email, password)
	if err != nil {
		_ = WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := generateToken(user.Email, user.Role, "AllYourBase")
	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// Set cookie with token.
	http.SetCookie(w, &http.Cookie{
		Name:     "clanplatform_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	})

	_ = WriteJSON(w, http.StatusOK, user)
}

func (api *api) AuthToken(w http.ResponseWriter, r *http.Request) {
	email, password, err := decodeCredentials(r)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.verifyUserCredentials(email, password)
	if err != nil {
		_ = WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := generateToken(user.Email, user.Role, "AllYourBase")
	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	_ = WriteJSON(w, http.StatusOK, map[string]string{"access_token": token})
}

func generateToken(email string, role entity.UserRoleEnum, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}
