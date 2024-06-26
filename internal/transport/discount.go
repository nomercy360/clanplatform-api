package transport

import (
	adm "clanplatform/internal/admin"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ListDiscountsHandler godoc
// @Summary List discounts
// @Tags discounts
// @Accept  json
// @Produce  json
// @Success 200 {array} Discount
// @Router /admin/discounts [get]
func (tr *transport) ListDiscountsHandler(c echo.Context) error {
	discounts, err := tr.admin.ListDiscounts()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, discounts)
}

// CreateDiscountHandler godoc
// @Summary Create discount
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param discount body CreateDiscount true "Discount data"
// @Success 200 {object} Discount
// @Router /admin/discounts [post]
func (tr *transport) CreateDiscountHandler(c echo.Context) error {
	var discount adm.CreateDiscount

	if err := c.Bind(&discount); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(discount); err != nil {
		return err
	}

	res, err := tr.admin.CreateDiscount(discount)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}
