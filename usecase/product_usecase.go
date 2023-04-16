package usecase

import (
	"context"
	"fmt"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/google/uuid"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Create(ctx context.Context, name string, userID uuid.UUID) (*domain.Product, error)
}

type ProductUsecase struct {
	productRepo ProductRepository
}

func NewProduct(productRepo ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
	}
}

func (u *ProductUsecase) View(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	pdt, err := u.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("productRepo.FindByID: %w", err)
	}

	if !pdt.IsPublished() {
		return nil, ErrProductNotFound
	}

	return pdt, nil
}

type CreateProductDto struct {
	Name   string
	UserID uuid.UUID
}

func (u *ProductUsecase) Create(ctx context.Context, dto CreateProductDto) (*domain.Product, error) {
	if name := domain.ProductName(dto.Name); !name.Valid() {
		return nil, ErrProductNameBadFormat
	}

	pdt, err := u.productRepo.Create(ctx, dto.Name, dto.UserID)
	if err != nil {
		return nil, fmt.Errorf("productRepo.Create: %w", err)
	}

	return pdt, nil
}

func (u *ProductUsecase) Delete(ctx context.Context, id, userID uuid.UUID) error {
	pdt, err := u.productRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("productRepo.FindByID: %w", err)
	}

	if !pdt.IsMine(userID) {
		return ErrProductUnauthorized
	}

	if err := u.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("productRepo.Delete: %w", err)
	}

	return nil
}
