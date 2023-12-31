package domain

import (
	"context"
)

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (svc *ProductService) PreparePurchase(ctx context.Context, unit int, p *Product, discounts []Discount) (*Purchase, error) {
	basePrice := p.Price

	var validDiscounts []Discount
	for _, d := range discounts {
		if !d.IsValid() {
			return nil, ErrNegativePrice
		}

		if unit < d.MinPurchaseQty {
			continue
		}

		validDiscounts = append(validDiscounts, d)
	}

	p, err := p.WithDiscount(validDiscounts...)
	if err != nil {
		return nil, err
	}

	return &Purchase{
		ProductID: p.ID,
		BasePrice: basePrice,
		Discount:  p.Price - basePrice,
		Unit:      unit,
	}, nil
}
