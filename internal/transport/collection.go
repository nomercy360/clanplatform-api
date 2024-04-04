package transport

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (tr *transport) ListCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	collections, err := tr.admin.ListCollections()

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, collections)
}

func (tr *transport) CreateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title  string `json:"title"`
		Handle string `json:"handle"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	collection, err := tr.admin.CreateCollection(data.Title, data.Handle)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusCreated, collection)
}

func (tr *transport) GetCollectionByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	collection, err := tr.admin.GetCollectionByID(id)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, collection)
}

func (tr *transport) UpdateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var data struct {
		Title  *string `json:"title"`
		Handle *string `json:"handle"`
	}

	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := tr.admin.UpdateCollection(data.Title, data.Handle, id)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, res)
}

func (tr *transport) DeleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = tr.admin.DeleteCollection(id)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, nil)
}

func (tr *transport) AddProductsToCollectionHandler(w http.ResponseWriter, r *http.Request) {
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

	err = tr.admin.AddProductsToCollection(collectionID, data.ProductIDs)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (tr *transport) RemoveProductsFromCollectionHandler(w http.ResponseWriter, r *http.Request) {
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

	err = tr.admin.RemoveProductsFromCollection(collectionID, data.ProductIDs)

	if err != nil {
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
