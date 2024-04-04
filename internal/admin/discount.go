package admin

import (
	"clanplatform/internal/db"
	"time"
)

func (adm *admin) ListDiscounts() ([]db.Discount, error) {
	disc, err := adm.storage.ListDiscounts()

	if err != nil {
		return nil, err
	}

	return disc, nil
}

func (adm *admin) CreateDiscount(discount db.Discount) (*db.Discount, error) {
	if discount.Code == "" {
		return nil, invalidReqErr
	}

	if discount.Value <= 0 {
		return nil, invalidReqErr
	}

	if discount.Type != db.Percentage && discount.Type != db.Fixed {
		return nil, invalidReqErr
	}

	if discount.StartsAt.IsZero() {
		discount.StartsAt = time.Now()
	}

	res, err := adm.storage.CreateDiscount(discount)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (adm *admin) UpdateDiscount(discount db.Discount) (*db.Discount, error) {
	if discount.ID == 0 {
		return nil, invalidReqErr
	}

	if discount.Code == "" {
		return nil, invalidReqErr
	}

	if discount.Value <= 0 {
		return nil, invalidReqErr
	}

	if discount.Type != db.Percentage && discount.Type != db.Fixed {
		return nil, invalidReqErr
	}

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
