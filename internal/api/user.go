package api

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
