package admin

import (
	"clanplatform/internal/db"
	"clanplatform/internal/terrors"
)

func (adm *admin) CreateProductVariant(variant db.ProductVariant) (*db.ProductVariant, error) {
	res, err := adm.storage.CreateProductVariant(variant.ProductID, variant.Title, variant.Inventory)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (adm *admin) ListProducts() ([]db.ProductWithDetails, error) {
	products, err := adm.storage.ListProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (adm *admin) GetProductByID(id int64) (*db.ProductWithDetails, error) {
	product, err := adm.storage.GetProductByID(id)

	if err != nil {
		if db.IsNoRowsError(err) {
			return nil, terrors.NotFound(err)
		}
		return nil, err
	}

	return product, nil
}

func (adm *admin) CreateProduct(product db.Product) (*db.Product, error) {
	if err := product.Validate(); err != nil {
		return nil, err
	}

	res, err := adm.storage.CreateProduct(product)

	if err != nil {
		return nil, err
	}

	return res, nil
}

type UpdateProductRequest struct {
	Title       *string `json:"title"`
	Subtitle    *string `json:"subtitle"`
	Description *string `json:"description"`
	Handle      *string `json:"handle"`
}

func (adm *admin) UpdateProduct(id int64, upd UpdateProductRequest) (*db.Product, error) {
	if id == 0 {
		return nil, invalidReqErr
	}

	update := make(map[string]interface{})
	if upd.Title != nil {
		update["title"] = *upd.Title
	}

	if upd.Subtitle != nil {
		update["subtitle"] = *upd.Subtitle
	}

	if upd.Description != nil {
		update["description"] = *upd.Description
	}

	if upd.Handle != nil {
		update["handle"] = *upd.Handle
	}

	res, err := adm.storage.UpdateProduct(id, update)

	if err != nil {
		return nil, err
	}

	return res, nil
}
