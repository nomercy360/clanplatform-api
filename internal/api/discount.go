package api

import (
	"clanplatform/internal/entity"
	"encoding/json"
	"net/http"
	"time"
)

func (api *api) ListDiscounts(w http.ResponseWriter, r *http.Request) {
	discounts, err := api.storage.GetDiscounts()

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, discounts)
}

func (api *api) CreateDiscount(w http.ResponseWriter, r *http.Request) {
	var discount entity.Discount
	if err := json.NewDecoder(r.Body).Decode(&discount); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the discount.
	if discount.Code == "" {
		_ = WriteError(w, http.StatusBadRequest, "code is required")
		return
	}

	if discount.Value <= 0 {
		_ = WriteError(w, http.StatusBadRequest, "value must be greater than 0")
		return
	}

	// validate discount type should be either "percentage" or "fixed"
	if discount.Type != entity.Percentage && discount.Type != entity.Fixed {
		_ = WriteError(w, http.StatusBadRequest, "type must be either 'percentage' or 'fixed'")
		return
	}

	if discount.StartsAt.IsZero() {
		discount.StartsAt = time.Now()
	}

	res, err := api.storage.CreateDiscount(discount)
	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusCreated, res)
}
