package transport

import (
	adm "clanplatform/internal/admin"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ListUsersHandler godoc
// @Summary List users
// @Description get users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} User
// @Router /admin/users [get]
func (tr *transport) ListUsersHandler(c echo.Context) error {
	users, err := tr.admin.ListUsers()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

// CreateUserHandler godoc
// @Summary Create user
// @Description create user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body CreateUser true "User data"
// @Success 201 {object} User
// @Router /admin/users [post]
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
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

// AuthCookieHandler godoc
// @Summary Authenticate user
// @Description authenticate user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body AuthUser true "User data"
// @Success 200 {object} UserWithToken
// @Router /admin/auth [post]
func (tr *transport) AuthCookieHandler(c echo.Context) error {
	var data adm.AuthUser

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	user, err := tr.admin.AuthUser(data)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cookie := new(http.Cookie)

	cookie.Name = "token"
	cookie.Value = user.Token
	cookie.HttpOnly = true
	cookie.Path = "/"

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, user)
}

// AuthTokenHandler godoc
// @Summary Authenticate user
// @Description authenticate user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body AuthUser true "User data"
// @Success 200 {object} UserWithToken
// @Router /admin/auth/token [post]
func (tr *transport) AuthTokenHandler(c echo.Context) error {
	var data adm.AuthUser

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	user, err := tr.admin.AuthUser(data)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}
