package db

import (
	"clanplatform/internal/entity"
)

func (s *storage) ListProducts() ([]entity.Product, error) {
	products := make([]entity.Product, 0)
	err := s.pg.Select(&products, "SELECT * FROM products")

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *storage) CreateProduct(title string, subtitle string, description string, handle string) (entity.Product, error) {
	query := `
		INSERT INTO products (title, subtitle, description, handle, is_published) 
		VALUES (:title, :subtitle, :description, :handle, :is_published)
		RETURNING *;
	`

	var result entity.Product

	rows, err := s.pg.NamedQuery(query, entity.Product{
		Title:       title,
		Subtitle:    subtitle,
		Description: description,
		Handle:      handle,
		IsPublished: false,
	})

	if err != nil {
		return result, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&result)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func (s *storage) CreateProductVariant(
	productID int64,
	title string,
	priceList []entity.Price,
	inventory int,
) (entity.ProductVariant, error) {
	query := `
		INSERT INTO product_variants (product_id, title, inventory)
		VALUES (:product_id, :title, :inventory)
		RETURNING *;
	`

	var result entity.ProductVariant

	rows, err := s.pg.NamedQuery(query, entity.ProductVariant{
		ProductID: productID,
		Title:     title,
		Inventory: inventory,
	})

	if err != nil {
		return result, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&result)
		if err != nil {
			return result, err
		}
	}

	for _, price := range priceList {
		_, err = s.pg.Exec(`
			INSERT INTO product_prices (variant_id, currency_id, amount)
			VALUES ($1, $2, $3)
		`, result.ID, price.CurrencyCode, price.Amount)

		if err != nil {
			return result, err
		}
	}

	return result, nil
}
