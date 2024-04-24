package transport

import (
	admSvc "clanplatform/internal/admin"
	"clanplatform/internal/db"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (tr *transport) CreateProductVariantHandler(c echo.Context) error {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var productVariant db.ProductVariant

	if err := c.Bind(&productVariant); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	productVariant.ProductID = productID

	if err := c.Validate(productVariant); err != nil {
		return err
	}

	res, err := tr.admin.CreateProductVariant(productVariant)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (tr *transport) ListProductsHandler(c echo.Context) error {
	products, err := tr.admin.ListProducts()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, products)
}

func (tr *transport) GetProductByIDHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	product, err := tr.admin.GetProductByID(id)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, product)
}

func (tr *transport) CreateProductHandler(c echo.Context) error {
	var product db.Product
	if err := c.Bind(&product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := tr.admin.CreateProduct(product)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (tr *transport) UpdateProductHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	var upd admSvc.UpdateProductRequest

	if err := c.Bind(&upd); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := tr.admin.UpdateProduct(id, upd)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
