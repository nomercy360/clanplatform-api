package transport

import (
	"clanplatform/internal/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (tr *transport) ListDiscountsHandler(c echo.Context) error {
	discounts, err := tr.admin.ListDiscounts()

	if err != nil {
		return WriteError(c.Response(), http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, discounts)
}

func (tr *transport) CreateDiscountHandler(c echo.Context) error {
	var discount db.Discount

	if err := c.Bind(&discount); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := tr.admin.CreateDiscount(discount)

	if err != nil {
		return WriteError(c.Response(), http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
