package usecase

import (
	"context"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/google/uuid"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
}

type ProductUsecase struct {
	productRepo ProductRepository
}

func NewProduct(productRepo ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
	}
}

func (p *ProductUsecase) View(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	pdt, err := p.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !pdt.IsPublished() {
		return nil, ErrProductNotFound
	}

	return pdt, nil
}

func (p *ProductUsecase) Delete(ctx context.Context, id, userID uuid.UUID) error {
	pdt, err := p.productRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !pdt.IsMine(userID) {
		return ErrProductUnauthorized
	}

	return nil
}
