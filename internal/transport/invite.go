package transport

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (tr *transport) InviteUserHandler(c echo.Context) error {
	var data struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := tr.admin.InviteUser(data.Role, data.Email)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "invite sent"})
}

func (tr *transport) ListInvitesHandler(c echo.Context) error {
	invites, err := tr.admin.ListInvites()

	if err != nil {
		return WriteError(c.Response(), http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, invites)
}

func (tr *transport) AcceptInviteHandler(c echo.Context) error {
	var data struct {
		Token string `json:"token"`
		User  struct {
			Password string `json:"password"`
			FullName string
		}
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := tr.admin.AcceptInvite(data.Token, data.User.Password, data.User.FullName)

	if err != nil {
		return WriteError(c.Response(), http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "invite accepted"})
}
