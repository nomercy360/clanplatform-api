package store

import (
	"clanplatform/internal/db"
)

func (st *store) ListProducts() ([]db.ProductWithDetails, error) {
	products, err := st.storage.ListProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (st *store) GetProductByID(id int64) (*db.ProductWithDetails, error) {
	product, err := st.storage.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
