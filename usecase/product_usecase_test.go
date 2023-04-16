package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/alextanhongpin/go-domain-test/domain/factories"
	"github.com/alextanhongpin/go-domain-test/mocks"
	"github.com/alextanhongpin/go-domain-test/types"
	"github.com/alextanhongpin/go-domain-test/usecase"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	isContext = mock.MatchedBy(func(ctx context.Context) bool {
		return true
	})
)

func TestProductUsecase(t *testing.T) {
	type stub struct {
		findByID    *domain.Product
		findByIDErr error
	}

	type args struct {
		id uuid.UUID
	}

	p := factories.NewProduct()

	tests := []struct {
		name    string
		args    args
		stub    stub
		want    *domain.Product
		wantErr error
	}{
		{
			name: "success",
			args: args{id: p.ID},
			stub: stub{findByID: p},
			want: p,
		},
		{
			name: "not yet published",
			args: args{id: p.ID},
			stub: stub{
				findByID: func() *domain.Product {
					cp := *p
					cp.PublishedAt = nil
					return &cp
				}(),
			},
			wantErr: usecase.ErrProductNotFound,
		},
		{
			name: "published in the future",
			args: args{id: p.ID},
			stub: stub{
				findByID: func() *domain.Product {
					cp := *p
					cp.PublishedAt = types.Ptr(time.Now().Add(1 * time.Second))
					return &cp
				}(),
			},
			wantErr: usecase.ErrProductNotFound,
		},
		{
			name:    "failed",
			args:    args{id: p.ID},
			stub:    stub{findByIDErr: errors.New("db error")},
			wantErr: errors.New("db error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			productRepo := new(mocks.ProductRepository)
			productRepo.On("FindByID", isContext, tc.args.id).Return(tc.stub.findByID, tc.stub.findByIDErr).Once()

			uc := usecase.NewProduct(productRepo)
			ucPdt, err := uc.View(context.Background(), tc.args.id)
			assert.Equal(tc.wantErr, err)
			if diff := cmp.Diff(tc.want, ucPdt); diff != "" {
				t.Errorf("want (+), got (-): %v", diff)
			}
		})
	}
}
