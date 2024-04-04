package admin

import (
	"clanplatform/internal/db"
)

func (adm *admin) CreateProductVariant(variant db.ProductVariant) (*db.ProductVariant, error) {
	if variant.ProductID == 0 {
		return nil, invalidReqErr
	}

	if variant.Title == "" {
		return nil, invalidReqErr
	}

	if variant.Inventory == 0 {
		return nil, invalidReqErr
	}

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
		return nil, err
	}

	return product, nil
}
