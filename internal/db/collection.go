package db

import (
	"clanplatform/internal/entity"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

func (s *Storage) ListCollections() ([]entity.ProductCollection, error) {
	collections := make([]entity.ProductCollection, 0)
	err := s.pg.Select(&collections, "SELECT * FROM product_collections")

	if err != nil {
		return nil, err
	}

	return collections, nil
}

func (s *Storage) CreateCollection(title string, handle string) (entity.ProductCollection, error) {
	query := `
		INSERT INTO product_collections (title, handle)
		VALUES (:title, :handle)
		RETURNING *;
	`

	collection := entity.ProductCollection{
		Title:  title,
		Handle: handle,
	}

	rows, err := s.pg.NamedQuery(query, collection)

	if err != nil {
		return entity.ProductCollection{}, entity.ErrDatabase
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&collection)
		if err != nil {
			return entity.ProductCollection{}, entity.ErrDatabase
		}
	}

	return collection, nil
}

func (s *Storage) GetCollectionByID(id int64) (entity.ProductCollection, error) {
	collection := entity.ProductCollection{}
	err := s.pg.Get(&collection, "SELECT * FROM product_collections WHERE id = $1", id)

	if err != nil {
		return entity.ProductCollection{}, err
	}

	return collection, nil
}

func (s *Storage) UpdateCollection(title *string, handle *string, id int64) (entity.ProductCollection, error) {
	var res entity.ProductCollection

	var updates []string
	params := map[string]interface{}{
		"id": id,
	}

	if title != nil {
		updates = append(updates, "title = :title")
		params["title"] = *title
	}

	if handle != nil {
		updates = append(updates, "handle = :handle")
		params["handle"] = *handle
	}

	query := fmt.Sprintf(`
		UPDATE product_collections
		SET %s
		WHERE id = :id
		RETURNING *;
	`, strings.Join(updates, ", "))

	rows, err := s.pg.NamedQuery(query, params)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&res)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (s *Storage) DeleteCollection(id int64) error {
	_, err := s.pg.Exec("DELETE FROM product_collections WHERE id = $1", id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddProductsToCollection(collectionID int64, productIDs []int64) error {
	query, args, err := sqlx.In(`
		UPDATE products
		SET collection_id = ?
		WHERE id IN (?)
	`, collectionID, productIDs)

	if err != nil {
		return err
	}

	query = s.pg.Rebind(query)

	_, err = s.pg.Exec(query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveProductsFromCollection(collectionID int64, productIDs []int64) error {
	query, args, err := sqlx.In(`
		UPDATE products
		SET collection_id = NULL
		WHERE collection_id = ? AND id IN (?)
	`, collectionID, productIDs)

	if err != nil {
		return err
	}

	query = s.pg.Rebind(query)

	_, err = s.pg.Exec(query, args...)

	if err != nil {
		return err
	}

	return nil
}
