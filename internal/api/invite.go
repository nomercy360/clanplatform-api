package api

import (
	"bytes"
	"clanplatform/internal/entity"
	"clanplatform/internal/services"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"html/template"
	"net/http"
	"time"
)

func (api *api) InviteUser(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	role := entity.UserRoleEnum(data.Role)

	if role != entity.Merchant && role != entity.Admin {
		_ = WriteError(w, http.StatusBadRequest, "invalid role")
		return
	}

	mySigningKey := []byte(api.jwtSecret)

	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		Subject:   data.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(mySigningKey)

	err = api.storage.InviteUser(ss, data.Email, role)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type Context struct {
		InviteURL string
	}

	tmpl, err := template.ParseFiles("templates/otp.gohtml")

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var tpl bytes.Buffer
	// TODO: change to real URL when frontend is ready
	err = tmpl.Execute(&tpl, Context{InviteURL: "http://localhost:8080/invite/accept?token=" + ss})

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := &services.MailMessage{
		To:       data.Email,
		Subject:  "You have been invited to join the platform",
		From:     "hi@mxksim.dev",
		HtmlBody: tpl.String(),
	}

	if err = api.email.SendEmail(message); err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (api *api) ListInvites(w http.ResponseWriter, r *http.Request) {
	invites, err := api.storage.ListInvites()

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, invites)
}

func (api *api) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string `json:"token"`
		User  struct {
			Password  string `json:"password"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if data.Token == "" || data.User.Password == "" || data.User.FirstName == "" || data.User.LastName == "" {
		_ = WriteError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	token, err := jwt.Parse(data.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(api.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		_ = WriteError(w, http.StatusBadRequest, "invalid token")
		return
	}

	email, _ := token.Claims.GetSubject()

	if email == "" {
		_ = WriteError(w, http.StatusBadRequest, "invalid token")
		return
	}

	invite, _ := api.storage.GetInviteByEmail(email)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if invite == nil {
		_ = WriteError(w, http.StatusNotFound, "not found")
		return
	}

	hashedPassword, _ := hashPassword(data.User.Password)

	_, err = api.storage.CreateUser(entity.User{
		Email:        email,
		PasswordHash: hashedPassword,
		FirstName:    data.User.FirstName,
		LastName:     data.User.LastName,
		Role:         invite.Role,
	})

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, map[string]string{"message": "invite accepted"})
}
