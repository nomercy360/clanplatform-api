package transport

import (
	"encoding/json"
	"net/http"
)

func (tr *transport) InviteUserHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (tr *transport) ListInvitesHandler(w http.ResponseWriter, r *http.Request) {
	invites, err := tr.admin.ListInvites()

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, invites)
}

func (tr *transport) AcceptInviteHandler(w http.ResponseWriter, r *http.Request) {
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

	err := tr.admin.AcceptInvite(data.Token, data.User.Password, data.User.FirstName, data.User.LastName)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, map[string]string{"message": "invite accepted"})
}
