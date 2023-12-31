package domain

import "github.com/google/uuid"

type Purchase struct {
	ProductID uuid.UUID
	BasePrice int
	Discount  int
	Unit      int
}
