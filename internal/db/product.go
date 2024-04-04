package db

import (
	"encoding/json"
	"fmt"
	"time"
)

type ProductPrice struct {
	ID         int64      `db:"id" json:"id"`
	ProductID  int64      `db:"product_id" json:"product_id"`
	CurrencyID string     `db:"currency_id" json:"currency_id"`
	Amount     int        `db:"amount" json:"amount"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`
}

type ProductCollection struct {
	ID        int64      `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	Handle    string     `db:"handle" json:"handle"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type Product struct {
	ID           int64            `db:"id" json:"id"`
	Title        string           `db:"title" json:"title"`
	Subtitle     string           `db:"subtitle" json:"subtitle"`
	Description  string           `db:"description" json:"description"`
	Handle       string           `db:"handle" json:"handle"`
	IsPublished  bool             `db:"is_published" json:"is_published"`
	CollectionID *int64           `db:"collection_id" json:"collection_id"`
	CreatedAt    time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time        `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time       `db:"deleted_at" json:"deleted_at"`
	Metadata     *json.RawMessage `db:"metadata" json:"metadata"`
}

type ProductWithDetails struct {
	ID           int64            `db:"id" json:"id"`
	Title        string           `db:"title" json:"title"`
	Subtitle     string           `db:"subtitle" json:"subtitle"`
	Description  string           `db:"description" json:"description"`
	Handle       string           `db:"handle" json:"handle"`
	IsPublished  bool             `db:"is_published" json:"is_published"`
	CollectionID *int64           `db:"collection_id" json:"collection_id"`
	CreatedAt    time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time        `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time       `db:"deleted_at" json:"deleted_at"`
	Metadata     *json.RawMessage `db:"metadata" json:"metadata"`
	Images       []Image          `json:"images"`
	Variants     []ProductVariant `json:"variants"`
	Categories   []Category       `json:"categories"`
	Prices       []ProductPrice   `json:"prices"`
}

type Currency struct {
	Code   string `db:"code" json:"code"`
	Name   string `db:"name" json:"name"`
	Symbol string `db:"symbol" json:"symbol"`
}

type Category struct {
	ID        int64      `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type Image struct {
	ID        int64      `db:"id" json:"id"`
	URL       string     `db:"url" json:"url"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	IsMain    bool       `db:"is_main" json:"is_main"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type ProductVariant struct {
	ID        int64     `db:"id"`
	ProductID int64     `db:"product_id"`
	Title     string    `db:"title"`
	Inventory int       `db:"inventory"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Helper function to check and add unique images
func addUniqueImage(images []Image, newImage Image) []Image {
	for _, img := range images {
		if img.ID == newImage.ID {
			return images
		}
	}
	return append(images, newImage)
}

// Helper function to check and add unique prices
func addUniquePrice(prices []ProductPrice, newPrice ProductPrice) []ProductPrice {
	for _, price := range prices {
		if price.ID == newPrice.ID {
			return prices
		}
	}
	return append(prices, newPrice)
}

func addUniqueCategory(categories []Category, newCategory Category) []Category {
	for _, category := range categories {
		if category.ID == newCategory.ID {
			return categories
		}
	}
	return append(categories, newCategory)
}

func addUniqueVariant(variants []ProductVariant, newVariant ProductVariant) []ProductVariant {
	for _, variant := range variants {
		if variant.ID == newVariant.ID {
			return variants
		}
	}

	return append(variants, newVariant)
}

type Price struct {
	Amount       int    `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

func (s *storage) ListProducts() ([]ProductWithDetails, error) {
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
		WHERE p.is_published = true;
	`

	var products []ProductWithDetails
	productMap := make(map[int64]*ProductWithDetails)

	// Execute the query
	rows, err := s.pg.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p ProductWithDetails
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

		product, ok := productMap[p.ID]
		if !ok {
			product = &p
			productMap[p.ID] = product
		}

		if imageID != nil {
			i := Image{
				ID:        *imageID,
				URL:       *imageURL,
				IsMain:    *isMain,
				CreatedAt: *imageCreated,
				UpdatedAt: *imageUpdated,
				DeletedAt: imageDeleted,
			}

			product.Images = addUniqueImage(product.Images, i)
		}

		if priceID != nil {
			pp := ProductPrice{
				ID:         *priceID,
				CurrencyID: *currencyCode,
				Amount:     *amount,
				CreatedAt:  *priceCreated,
				UpdatedAt:  *priceUpdated,
				DeletedAt:  priceDeleted,
			}

			product.Prices = addUniquePrice(product.Prices, pp)
		}

		if variantID != nil {
			pv := ProductVariant{
				ID:        *variantID,
				Title:     *variantTitle,
				Inventory: *inventory,
				CreatedAt: *variantCreated,
				UpdatedAt: *variantUpdated,
			}

			product.Variants = addUniqueVariant(product.Variants, pv)
		}

		if categoryID != nil {
			category := Category{
				ID:        *categoryID,
				Name:      *categoryName,
				CreatedAt: *categoryCreated,
				UpdatedAt: *categoryUpdated,
				DeletedAt: categoryDeleted,
			}

			product.Categories = addUniqueCategory(product.Categories, category)
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

func (s *storage) CreateProduct(title string, subtitle string, description string, handle string) (Product, error) {
	query := `
		INSERT INTO products (title, subtitle, description, handle, is_published) 
		VALUES (:title, :subtitle, :description, :handle, :is_published)
		RETURNING *;
	`

	var result Product

	rows, err := s.pg.NamedQuery(query, Product{
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
	inventory int,
) (*ProductVariant, error) {
	query := `
		INSERT INTO product_variants (product_id, title, inventory)
		VALUES (:product_id, :title, :inventory)
		RETURNING *;
	`

	var result ProductVariant

	rows, err := s.pg.NamedQuery(query, ProductVariant{
		ProductID: productID,
		Title:     title,
		Inventory: inventory,
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (s *storage) GetProductByID(id int64) (*ProductWithDetails, error) {
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
		WHERE p.id = $1
		ORDER BY p.id;
	`

	var product ProductWithDetails

	rows, err := s.pg.Query(query, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
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
			&product.ID, &product.Title, &product.Subtitle, &product.Description, &product.Handle, &product.IsPublished,
			&product.CollectionID, &product.Metadata, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt,
			&imageID, &imageURL, &isMain, &imageCreated, &imageUpdated, &imageDeleted,
			&priceID, &currencyCode, &amount, &priceCreated, &priceUpdated, &priceDeleted,
			&variantID, &variantTitle, &inventory, &variantCreated, &variantUpdated, &variantDeleted,
			&categoryID, &categoryName, &categoryCreated, &categoryUpdated, &categoryDeleted,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning product row: %w", err)
		}

		if imageID != nil {
			i := Image{
				ID:        *imageID,
				URL:       *imageURL,
				IsMain:    *isMain,
				CreatedAt: *imageCreated,
				UpdatedAt: *imageUpdated,
				DeletedAt: imageDeleted,
			}

			product.Images = addUniqueImage(product.Images, i)
		}

		if priceID != nil {
			pp := ProductPrice{
				ID:         *priceID,
				CurrencyID: *currencyCode,
				Amount:     *amount,
				CreatedAt:  *priceCreated,
				UpdatedAt:  *priceUpdated,
				DeletedAt:  priceDeleted,
			}

			product.Prices = addUniquePrice(product.Prices, pp)
		}

		if variantID != nil {
			pv := ProductVariant{
				ID:        *variantID,
				Title:     *variantTitle,
				Inventory: *inventory,
				CreatedAt: *variantCreated,
				UpdatedAt: *variantUpdated,
			}

			product.Variants = addUniqueVariant(product.Variants, pv)
		}

		if categoryID != nil {
			category := Category{
				ID:        *categoryID,
				Name:      *categoryName,
				CreatedAt: *categoryCreated,
				UpdatedAt: *categoryUpdated,
				DeletedAt: categoryDeleted,
			}

			product.Categories = addUniqueCategory(product.Categories, category)
		}

	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating product rows: %w", err)
	}

	return &product, nil
}
