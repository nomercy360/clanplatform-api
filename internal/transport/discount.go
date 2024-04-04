package transport

import (
	"clanplatform/internal/db"
	"encoding/json"
	"net/http"
)

func (tr *transport) ListDiscountsHandler(w http.ResponseWriter, r *http.Request) {
	discounts, err := tr.admin.ListDiscounts()

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, discounts)
}

func (tr *transport) CreateDiscountHandler(w http.ResponseWriter, r *http.Request) {
	var discount db.Discount
	if err := json.NewDecoder(r.Body).Decode(&discount); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := tr.admin.CreateDiscount(discount)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusCreated, res)
}
