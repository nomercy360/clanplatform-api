package db

import (
	"clanplatform/internal/entity"
)

func (s *storage) CreateDiscount(discount entity.Discount) (entity.Discount, error) {
	query := `
		INSERT INTO discounts (code, is_active, type, usage_limit, ends_at, value, starts_at)
		VALUES (:code, :is_active, :type, :usage_limit, :ends_at, :value, :starts_at)
		RETURNING *;
	`

	rows, err := s.pg.NamedQuery(query, discount)

	if err != nil {
		return entity.Discount{}, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&discount)
		if err != nil {
			return entity.Discount{}, err
		}
	}

	return discount, nil
}

func (s *storage) GetDiscounts() ([]entity.Discount, error) {
	discounts := make([]entity.Discount, 0)
	err := s.pg.Select(&discounts, "SELECT * FROM discounts")
	if err != nil {
		return nil, err
	}
	return discounts, nil
}

func (s *storage) UpdateDiscount(discount entity.Discount) (entity.Discount, error) {
	query := `
		UPDATE discounts
		SET code = :code, is_active = :is_active, value = :value, type = :type, usage_limit = :usage_limit, starts_at = :starts_at, ends_at = :ends_at
		WHERE id = :id
		RETURNING *;
	`

	rows, err := s.pg.NamedQuery(query, discount)

	if err != nil {
		return entity.Discount{}, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&discount)
		if err != nil {
			return entity.Discount{}, err
		}
	}

	return discount, nil
}

func (s *storage) DeleteDiscount(id string) error {
	query := `
		DELETE FROM discounts
		WHERE id = $1;
	`

	res, err := s.pg.Exec(query, id)

	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}
