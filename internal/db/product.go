package db

import (
	"clanplatform/internal/entity"
	"fmt"
	"time"
)

// Helper function to check and add unique images
func addUniqueImage(images []entity.Image, newImage entity.Image) []entity.Image {
	for _, img := range images {
		if img.ID == newImage.ID {
			return images
		}
	}
	return append(images, newImage)
}

// Helper function to check and add unique prices
func addUniquePrice(prices []entity.ProductPrice, newPrice entity.ProductPrice) []entity.ProductPrice {
	for _, price := range prices {
		if price.ID == newPrice.ID {
			return prices
		}
	}
	return append(prices, newPrice)
}

func addUniqueCategory(categories []entity.Category, newCategory entity.Category) []entity.Category {
	for _, category := range categories {
		if category.ID == newCategory.ID {
			return categories
		}
	}
	return append(categories, newCategory)
}

func addUniqueVariant(variants []entity.ProductVariant, newVariant entity.ProductVariant) []entity.ProductVariant {
	for _, variant := range variants {
		if variant.ID == newVariant.ID {
			return variants
		}
	}
	return append(variants, newVariant)
}

func (s *storage) ListProducts() ([]entity.ProductWithDetails, error) {
	query := `
		SELECT p.id, p.title, p.subtitle, p.description, p.handle, p.is_published,
		       p.collection_id, p.metadata, p.created_at, p.updated_at, p.deleted_at,
		       i.id, i.url, i.is_main, i.created_at, i.updated_at, i.deleted_at,
		       pp.id, pp.currency_id, pp.amount, pp.created_at, pp.updated_at, pp.deleted_at,
		       pv.id, pv.title, pv.inventory, pv.created_at, pv.updated_at, pv.deleted_at,
		       pc.id, pc.name, pc.created_at, pc.updated_at, pc.deleted_at
		FROM products p
		LEFT JOIN product_images pi ON p.id = pi.product_id
		LEFT JOIN images i ON pi.image_id = i.id
		LEFT JOIN product_prices pp ON p.id = pp.product_id
		LEFT JOIN product_variants pv ON p.id = pv.product_id
		LEFT JOIN product_categories_products pcp ON p.id = pcp.product_id
		LEFT JOIN product_categories pc ON pcp.category_id = pc.id
		WHERE p.is_published = true
		ORDER BY p.id;
	`

	var products []entity.ProductWithDetails
	productMap := make(map[int64]*entity.ProductWithDetails)

	// Execute the query
	rows, err := s.pg.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p entity.ProductWithDetails
		var imageID, priceID, variantID, categoryID *int64
		var imageURL, variantTitle, categoryName *string
		var isMain *bool
		var amount, inventory *int
		var currencyCode *string
		var categoryCreated, categoryUpdated, categoryDeleted *time.Time
		var priceCreated, priceUpdated, priceDeleted *time.Time
		var variantCreated, variantUpdated, variantDeleted *time.Time
		var imageCreated, imageUpdated, imageDeleted *time.Time

		err := rows.Scan(
			&p.ID, &p.Title, &p.Subtitle, &p.Description, &p.Handle, &p.IsPublished,
			&p.CollectionID, &p.Metadata, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
			&imageID, &imageURL, &isMain, &imageCreated, &imageUpdated, &imageDeleted,
			&priceID, &currencyCode, &amount, &priceCreated, &priceUpdated, &priceDeleted,
			&variantID, &variantTitle, &inventory, &variantCreated, &variantUpdated, &variantDeleted,
			&categoryID, &categoryName, &categoryCreated, &categoryUpdated, &categoryDeleted,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning product row: %w", err)
		}

		if imageID != nil {
			i := entity.Image{
				ID:        *imageID,
				URL:       *imageURL,
				IsMain:    *isMain,
				CreatedAt: *imageCreated,
				UpdatedAt: *imageUpdated,
				DeletedAt: imageDeleted,
			}

			p.Images = append(p.Images, i)
		}

		if priceID != nil {
			pp := entity.ProductPrice{
				ID:         *priceID,
				CurrencyID: *currencyCode,
				Amount:     *amount,
				CreatedAt:  *priceCreated,
				UpdatedAt:  *priceUpdated,
				DeletedAt:  priceDeleted,
			}

			p.Prices = append(p.Prices, pp)
		}

		if variantID != nil {
			pv := entity.ProductVariant{
				ID:        *variantID,
				Title:     *variantTitle,
				Inventory: *inventory,
				CreatedAt: *variantCreated,
				UpdatedAt: *variantUpdated,
			}

			p.Variants = append(p.Variants, pv)
		}

		if categoryID != nil {
			category := entity.Category{
				ID:        *categoryID,
				Name:      *categoryName,
				CreatedAt: *categoryCreated,
				UpdatedAt: *categoryUpdated,
				DeletedAt: categoryDeleted,
			}

			p.Categories = append(p.Categories, category)
		}

		if product, ok := productMap[p.ID]; ok {
			for _, img := range p.Images {
				product.Images = addUniqueImage(product.Images, img)
			}

			for _, price := range p.Prices {
				product.Prices = addUniquePrice(product.Prices, price)
			}

			for _, variant := range p.Variants {
				product.Variants = addUniqueVariant(product.Variants, variant)
			}

			for _, category := range p.Categories {
				product.Categories = addUniqueCategory(product.Categories, category)
			}

			productMap[p.ID] = product
		} else {
			productMap[p.ID] = &p
		}
	}

	for _, product := range productMap {
		products = append(products, *product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating product rows: %w", err)
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
			INSERT INTO product_prices (product_id, currency_id, amount)
			VALUES ($1, $2, $3)
		`, result.ID, price.CurrencyCode, price.Amount)

		if err != nil {
			return result, err
		}
	}

	return result, nil
}
