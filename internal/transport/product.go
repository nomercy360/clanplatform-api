package transport

import (
	"clanplatform/internal/db"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (tr *transport) CreateProductVariantHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title     string `json:"title"`
		Inventory int    `json:"inventory"`
	}

	productID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	productVariant := db.ProductVariant{
		ProductID: productID,
		Title:     request.Title,
		Inventory: request.Inventory,
	}

	res, err := tr.admin.CreateProductVariant(productVariant)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusCreated, res)
}

func (tr *transport) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := tr.admin.ListProducts()
	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, products)
}

func (tr *transport) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := tr.admin.GetProductByID(id)
	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, product)
}
