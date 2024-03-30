package api

import (
	"clanplatform/internal/entity"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (api *api) CreateProductVariant(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title     string         `json:"title"`
		PriceList []entity.Price `json:"price_list"`
		Inventory int            `json:"inventory"`
	}

	productID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if productID == 0 {
		_ = WriteError(w, http.StatusBadRequest, "product_id is required")
		return
	}

	if request.Title == "" {
		_ = WriteError(w, http.StatusBadRequest, "title is required")
		return
	}

	if len(request.PriceList) == 0 {
		_ = WriteError(w, http.StatusBadRequest, "price_list is required")
		return
	}

	if request.Inventory == 0 {
		_ = WriteError(w, http.StatusBadRequest, "inventory is required")
		return
	}

	productVariant, err := api.storage.CreateProductVariant(productID, request.Title, request.PriceList, request.Inventory)
	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusCreated, productVariant)
}

func (api *api) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := api.storage.ListProducts()
	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, products)
}
