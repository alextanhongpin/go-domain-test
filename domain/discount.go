package domain

import "github.com/google/uuid"

type Discount struct {
	ID             int64
	Name           string
	ProductID      uuid.UUID
	Amount         int // -tive amount for price deduction.
	MinPurchaseQty int
}

func (d *Discount) IsValid() bool {
	return d.Amount < 0 && d.MinPurchaseQty > 0
}
