package usecase

import (
	"context"
	"fmt"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/google/uuid"
)

type purchaseRepository interface {
	CheckUserEligibility(ctx context.Context, userID uuid.UUID) error
	FindProduct(ctx context.Context, productID uuid.UUID) (*domain.Product, error)
	FindProductDiscount(ctx context.Context, productID uuid.UUID) ([]domain.Discount, error)
	CreatePurchase(ctx context.Context, purchase domain.Purchase) error
}

type PurchaseUsecase struct {
	repo purchaseRepository
	svc  *domain.ProductService
}

func NewPurchaseUsecase(repo purchaseRepository) *PurchaseUsecase {
	return &PurchaseUsecase{
		repo: repo,
		svc:  domain.NewProductService(),
	}
}

type PurchaseDto struct {
	ProductID uuid.UUID
	UserID    uuid.UUID
	Unit      int
}

func (u *PurchaseUsecase) Purchase(ctx context.Context, dto PurchaseDto) error {
	if err := u.repo.CheckUserEligibility(ctx, dto.UserID); err != nil {
		return err
	}

	p, err := u.repo.FindProduct(ctx, dto.ProductID)
	if err != nil {
		return err
	}
	if !p.IsPublished() {
		return ErrProductNotFound
	}

	ds, err := u.repo.FindProductDiscount(ctx, dto.ProductID)
	if err != nil {
		return err
	}

	req, err := u.svc.PreparePurchase(ctx, dto.Unit, p, ds)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDiscountInvalid, err)
	}

	return u.repo.CreatePurchase(ctx, *req)
}
