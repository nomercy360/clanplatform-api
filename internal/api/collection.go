package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (api *api) ListCollections(w http.ResponseWriter, r *http.Request) {
	collections, err := api.storage.ListCollections()

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, collections)
}

func (api *api) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title  string `json:"title"`
		Handle string `json:"handle"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if data.Title == "" || data.Handle == "" {
		_ = WriteError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	collection, err := api.storage.CreateCollection(data.Title, data.Handle)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusCreated, collection)
}

func (api *api) GetCollection(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	collection, err := api.storage.GetCollectionByID(id)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, collection)
}

func (api *api) UpdateCollection(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var data struct {
		Title  *string `json:"title"`
		Handle *string `json:"handle"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	collection, err := api.storage.UpdateCollection(data.Title, data.Handle, id)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, collection)
}

func (api *api) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.storage.DeleteCollection(id)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, nil)
}

func (api *api) AddProductsToCollection(w http.ResponseWriter, r *http.Request) {
	collectionID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var data struct {
		ProductIDs []int64 `json:"product_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(data.ProductIDs) == 0 {
		_ = WriteError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	err = api.storage.AddProductsToCollection(collectionID, data.ProductIDs)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *api) RemoveProductsFromCollection(w http.ResponseWriter, r *http.Request) {
	collectionID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var data struct {
		ProductIDs []int64 `json:"product_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(data.ProductIDs) == 0 {
		_ = WriteError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	err = api.storage.RemoveProductsFromCollection(collectionID, data.ProductIDs)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
