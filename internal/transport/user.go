package transport

import (
	"encoding/json"
	"net/http"
	"time"
)

func (tr *transport) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := tr.admin.ListUsers()

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

// AuthCookieHandler set cookie and returns user
func (tr *transport) AuthCookieHandler(w http.ResponseWriter, r *http.Request) {
	email, password, err := decodeCredentials(r)

	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := tr.admin.AuthUser(email, password)

	http.SetCookie(w, &http.Cookie{
		Name:     "clanplatform_token",
		Value:    user.Token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	})

	_ = WriteJSON(w, http.StatusOK, user)
}

// AuthTokenHandler returns jwt token
func (tr *transport) AuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	email, password, err := decodeCredentials(r)

	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := tr.admin.AuthUser(email, password)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, map[string]string{"access_token": user.Token})
}
