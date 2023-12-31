package factories

import (
	"log"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/google/uuid"
)

func NewDiscount(variants ...string) *domain.Discount {
	// Valid discount by default.
	dis := &domain.Discount{
		ID:             1,
		Name:           "5$ off if you buy 2",
		ProductID:      uuid.New(),
		Amount:         -5,
		MinPurchaseQty: 2,
	}

	for _, v := range variants {
		switch v {
		case "positive_amount":
			dis.Amount = 5
		case "negative_purchase_qty":
			dis.MinPurchaseQty = -1
		default:
			log.Fatalf("unknown Discount variant: %s", v)
		}
	}

	return dis
}
