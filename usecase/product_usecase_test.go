package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/alextanhongpin/go-core-microservice/types/errors"
	"github.com/alextanhongpin/go-domain-test/domain/factories"
	"github.com/alextanhongpin/go-domain-test/mocks"
	"github.com/alextanhongpin/go-domain-test/types"
	"github.com/alextanhongpin/go-domain-test/usecase"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var isContext = mock.MatchedBy(func(ctx context.Context) bool {
	return true
})

func TestProductUsecase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		assert := assert.New(t)

		pdt := factories.NewProduct()
		productRepo := new(mocks.ProductRepository)
		productRepo.On("FindByID", isContext, pdt.ID).Return(pdt, nil).Once()

		uc := usecase.NewProduct(productRepo)
		ucPdt, err := uc.View(context.Background(), pdt.ID)
		assert.Nil(err)
		if diff := cmp.Diff(pdt, ucPdt); diff != "" {
			t.Errorf("want (+), got (-): %v", diff)
		}
	})

	t.Run("not published", func(t *testing.T) {
		assert := assert.New(t)

		pdt := factories.NewProduct()
		pdt.PublishedAt = types.Ptr(time.Now().Add(1 * time.Second)) // Only published 1s in the future.

		productRepo := new(mocks.ProductRepository)
		productRepo.On("FindByID", isContext, pdt.ID).Return(pdt, nil).Once()

		uc := usecase.NewProduct(productRepo)
		ucPdt, err := uc.View(context.Background(), pdt.ID)
		assert.Equal(usecase.ErrProductNotFound, err)
		assert.Nil(ucPdt)
	})

	t.Run("failed", func(t *testing.T) {
		assert := assert.New(t)

		productRepo := new(mocks.ProductRepository)
		productRepo.On("FindByID", isContext, uuid.Nil).Return(nil, errors.New("db error")).Once()

		uc := usecase.NewProduct(productRepo)
		ucPdt, err := uc.View(context.Background(), uuid.Nil)
		assert.Equal(errors.New("db error"), err)
		assert.Nil(ucPdt)
	})
}
