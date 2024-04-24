package admin

import (
	"clanplatform/internal/db"
	"clanplatform/internal/terrors"
	"time"
)

func (adm *admin) ListDiscounts() ([]db.Discount, error) {
	disc, err := adm.storage.ListDiscounts()

	if err != nil {
		return nil, err
	}

	return disc, nil
}

type CreateDiscount struct {
	Code       string    `json:"code" validate:"required"`
	Type       string    `json:"type" validate:"required,oneof=percentage fixed free_shipping"`
	Value      int       `json:"value" validate:"required,min=1"`
	StartsAt   time.Time `json:"starts_at"`
	UsageLimit int       `json:"usage_limit" validate:"omitempty,min=1"`
}

func (cd CreateDiscount) toDiscount() db.Discount {
	return db.Discount{
		Code:       cd.Code,
		Type:       cd.Type,
		Value:      cd.Value,
		StartsAt:   cd.StartsAt,
		UsageLimit: cd.UsageLimit,
	}
}

func (adm *admin) CreateDiscount(cd CreateDiscount) (*db.Discount, error) {
	if cd.StartsAt.IsZero() {
		cd.StartsAt = time.Now()
	}

	res, err := adm.storage.CreateDiscount(cd.toDiscount())

	if err != nil {
		if db.IsDuplicationError(err) {
			return nil, terrors.BadRequest(err)
		}

		return nil, err
	}

	return res, nil
}

func (adm *admin) UpdateDiscount(discount db.Discount) (*db.Discount, error) {
	res, err := adm.storage.UpdateDiscount(discount)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (adm *admin) DeleteDiscount(id string) error {
	if id == "" {
		return invalidReqErr
	}

	if err := adm.storage.DeleteDiscount(id); err != nil {
		return err
	}

	return nil
}
