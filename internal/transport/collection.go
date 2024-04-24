package transport

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (tr *transport) ListCollectionsHandler(c echo.Context) error {
	collections, err := tr.admin.ListCollections()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, collections)
}

func (tr *transport) CreateCollectionHandler(c echo.Context) error {
	var data struct {
		Title  string `json:"title"`
		Handle string `json:"handle"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	collection, err := tr.admin.CreateCollection(data.Title, data.Handle)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, collection)
}

func (tr *transport) GetCollectionByIDHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	collection, err := tr.admin.GetCollectionByID(id)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, collection)
}

func (tr *transport) UpdateCollectionHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var data struct {
		Title  *string `json:"title"`
		Handle *string `json:"handle"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	res, err := tr.admin.UpdateCollection(data.Title, data.Handle, id)

	if err != nil {
		return WriteError(c.Response(), http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (tr *transport) DeleteCollectionHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = tr.admin.DeleteCollection(id)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "collection deleted"})
}

func (tr *transport) AddProductsToCollectionHandler(c echo.Context) error {
	collectionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var data struct {
		ProductIDs []int64 `json:"product_ids"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	err = tr.admin.AddProductsToCollection(collectionID, data.ProductIDs)

	if err != nil {
		return WriteError(c.Response(), http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (tr *transport) RemoveProductsFromCollectionHandler(c echo.Context) error {
	collectionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return WriteError(c.Response(), http.StatusBadRequest, err.Error())
	}

	var data struct {
		ProductIDs []int64 `json:"product_ids"`
	}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = tr.admin.RemoveProductsFromCollection(collectionID, data.ProductIDs)

	if err != nil {
		return WriteError(c.Response(), http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
