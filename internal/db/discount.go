package db

import "time"

type Discount struct {
	ID         int64      `db:"id" json:"id"`
	Code       string     `db:"code" json:"code"`
	IsActive   bool       `db:"is_active" json:"is_active"`
	Type       string     `db:"type" json:"type"`
	UsageLimit int        `db:"usage_limit" json:"usage_limit"`
	UsageCount int        `db:"usage_count" json:"usage_count"`
	StartsAt   time.Time  `db:"starts_at" json:"starts_at"`
	EndsAt     *time.Time `db:"ends_at" json:"ends_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`
	Value      int        `db:"value" json:"value"`
}

func (s *Storage) CreateDiscount(discount Discount) (*Discount, error) {
	var res Discount

	query := `
		INSERT INTO discounts (code, is_active, type, usage_limit, ends_at, value, starts_at)
		VALUES (:code, :is_active, :type, :usage_limit, :ends_at, :value, :starts_at)
		RETURNING id, code, is_active, type, usage_limit, usage_count, starts_at, ends_at, created_at, updated_at, deleted_at, value;
	`

	rows, err := s.pg.NamedQuery(query, discount)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *Storage) ListDiscounts() ([]Discount, error) {
	discounts := make([]Discount, 0)
	err := s.pg.Select(&discounts, "SELECT * FROM discounts")
	if err != nil {
		return nil, err
	}
	return discounts, nil
}

func (s *Storage) UpdateDiscount(discount Discount) (*Discount, error) {
	var res Discount

	query := `
		UPDATE discounts
		SET code = :code, is_active = :is_active, value = :value, type = :type, usage_limit = :usage_limit, starts_at = :starts_at, ends_at = :ends_at
		WHERE id = :id
		RETURNING *;
	`

	rows, err := s.pg.NamedQuery(query, discount)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *Storage) DeleteDiscount(id string) error {
	query := `
		DELETE FROM discounts
		WHERE id = $1;
	`

	res, err := s.pg.Exec(query, id)

	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
