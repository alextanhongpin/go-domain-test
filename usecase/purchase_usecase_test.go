package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/alextanhongpin/go-domain-test/domain/factories"
	mocks "github.com/alextanhongpin/go-domain-test/mocks/github.com/alextanhongpin/go-domain-test/usecase"
	"github.com/alextanhongpin/go-domain-test/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPurchaseFlow(t *testing.T) {
	var wantErr = errors.New("want error")
	t.Run("success", func(t *testing.T) {
		f := newPurchaseFlow()
		assert.Nil(t, f.exec())
	})

	t.Run("check user eligibility error", func(t *testing.T) {
		f := newPurchaseFlow()
		f.stub.checkUserEligibility.err = wantErr
		assert.ErrorIs(t, f.exec(), wantErr)
	})

	t.Run("find product error", func(t *testing.T) {
		f := newPurchaseFlow()
		f.stub.findProduct.err = wantErr
		assert.ErrorIs(t, f.exec(), wantErr)
	})

	t.Run("product not published", func(t *testing.T) {
		f := newPurchaseFlow()
		f.stub.findProduct.data = factories.NewProduct("no_published_at")
		assert.ErrorIs(t, f.exec(), usecase.ErrProductNotFound)
	})

	t.Run("find product discount error", func(t *testing.T) {
		f := newPurchaseFlow()
		f.stub.findProductDiscount.err = wantErr
		assert.ErrorIs(t, f.exec(), wantErr)
	})

	t.Run("discount is invalid", func(t *testing.T) {
		f := newPurchaseFlow()
		f.stub.findProductDiscount.data = append(f.stub.findProductDiscount.data, *factories.NewDiscount("positive_amount"))
		assert.ErrorIs(t, f.exec(), usecase.ErrDiscountInvalid)
	})

	t.Run("create purchase error", func(t *testing.T) {
		f := newPurchaseFlow()
		f.stub.createPurchase.err = wantErr
		assert.ErrorIs(t, f.exec(), wantErr)
	})
}

type purchaseFlow struct {
	args usecase.PurchaseDto
	stub struct {
		checkUserEligibility arg0[uuid.UUID]
		findProduct          arg1[uuid.UUID, *domain.Product]
		findProductDiscount  arg1[uuid.UUID, []domain.Discount]
		createPurchase       arg0[domain.Purchase]
	}
}

func newPurchaseFlow() *purchaseFlow {
	f := new(purchaseFlow)

	p := factories.NewProduct()
	d := factories.NewDiscount()
	d.ProductID = p.ID

	f.args = usecase.PurchaseDto{
		ProductID: p.ID,
		UserID:    uuid.New(),
		Unit:      2,
	}

	f.stub.checkUserEligibility.args = f.args.UserID

	f.stub.findProduct.args = f.args.ProductID
	f.stub.findProduct.data = p

	f.stub.findProductDiscount.args = f.args.ProductID
	f.stub.findProductDiscount.data = []domain.Discount{*d}

	// For fields that needs to be updated after reload.
	if err := f.reload(); err != nil {
		panic(err)
	}

	return f
}

func (f *purchaseFlow) reload() error {
	svc := domain.NewProductService()
	req, err := svc.PreparePurchase(context.Background(), f.args.Unit, f.stub.findProduct.data, f.stub.findProductDiscount.data)
	if err != nil {
		return err
	}
	f.stub.createPurchase.args = *req

	return nil
}

func (f *purchaseFlow) exec() error {
	args := f.args
	stub := f.stub

	ctx := context.Background()

	repo := new(mocks.MockPurchaseRepository)
	repo.EXPECT().CheckUserEligibility(ctx, stub.checkUserEligibility.args).Return(stub.checkUserEligibility.err)
	repo.EXPECT().FindProduct(ctx, stub.findProduct.args).Return(stub.findProduct.data, stub.findProduct.err)
	repo.EXPECT().FindProductDiscount(ctx, stub.findProductDiscount.args).Return(stub.findProductDiscount.data, stub.findProductDiscount.err)
	repo.EXPECT().CreatePurchase(ctx, stub.createPurchase.args).Return(stub.createPurchase.err)

	u := usecase.NewPurchaseUsecase(repo)
	return u.Purchase(ctx, args)
}
