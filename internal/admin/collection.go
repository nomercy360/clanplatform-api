package admin

import (
	"clanplatform/internal/db"
	"errors"
)

var (
	invalidReqErr = errors.New("invalid request")
)

func (adm *admin) CreateCollection(title string, handle string) (*db.ProductCollection, error) {
	if title == "" || handle == "" {
		return nil, invalidReqErr
	}

	collection, err := adm.storage.CreateCollection(title, handle)

	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (adm *admin) GetCollectionByID(id int64) (*db.ProductCollection, error) {
	collection, err := adm.storage.GetCollectionByID(id)

	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (adm *admin) AddProductsToCollection(collectionID int64, productIDs []int64) error {
	if len(productIDs) == 0 {
		return invalidReqErr
	}

	if err := adm.storage.AddProductsToCollection(collectionID, productIDs); err != nil {
		return err
	}

	return nil
}

func (adm *admin) RemoveProductsFromCollection(collectionID int64, productIDs []int64) error {
	if len(productIDs) == 0 {
		return invalidReqErr
	}

	if err := adm.storage.RemoveProductsFromCollection(collectionID, productIDs); err != nil {
		return err
	}

	return nil
}

func (adm *admin) DeleteCollection(collectionID int64) error {
	if err := adm.storage.DeleteCollection(collectionID); err != nil {
		return err
	}

	return nil
}

func (adm *admin) ListCollections() ([]db.ProductCollection, error) {
	collections, err := adm.storage.ListCollections()

	if err != nil {
		return nil, err
	}

	return collections, nil
}

func (adm *admin) UpdateCollection(title *string, handle *string, id int64) (*db.ProductCollection, error) {
	if title == nil && handle == nil {
		return nil, invalidReqErr
	}

	collection, err := adm.storage.GetCollectionByID(id)
	if err != nil {
		return nil, err
	}

	if title != nil {
		collection.Title = *title
	}

	if handle != nil {
		collection.Handle = *handle
	}

	if _, err := adm.storage.UpdateCollection(&collection.Title, &collection.Handle, id); err != nil {
		return nil, err
	}

	return collection, nil
}
