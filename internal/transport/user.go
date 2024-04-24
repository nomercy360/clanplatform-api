package transport

import (
	adm "clanplatform/internal/admin"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (tr *transport) ListUsersHandler(c echo.Context) error {
	users, err := tr.admin.ListUsers()

	if err != nil {
		return WriteError(c.Response(), http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (tr *transport) CreateUserHandler(c echo.Context) error {
	var data adm.CreateUser

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	res, err := tr.admin.CreateUser(data)

	if err != nil {
		return WriteError(c.Response(), http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

// AuthCookieHandler set cookie and returns user
func (tr *transport) AuthCookieHandler(c echo.Context) error {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := tr.admin.AuthUser(data.Email, data.Password)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cookie := new(http.Cookie)

	cookie.Name = "token"
	cookie.Value = user.Token
	cookie.Expires = time.Now().Add(24 * time.Hour)

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, user)
}

// AuthTokenHandler returns jwt token
func (tr *transport) AuthTokenHandler(c echo.Context) error {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := tr.admin.AuthUser(data.Email, data.Password)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"token": user.Token})
}
