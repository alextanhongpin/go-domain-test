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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	isContext = mock.MatchedBy(func(ctx context.Context) bool {
		return true
	})

	// unwrapErr attempts to unwrap the wrapped error, otherwise it returns the
	// original error if there unwrapped error is nil.
	unwrapErr = func(err error) error {
		e := errors.Unwrap(err)
		if e != nil {
			return e
		}

		return err
	}
)

func TestProductUsecaseView(t *testing.T) {
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
			args, stub := tc.args, tc.stub

			productRepo := new(mocks.ProductRepository)
			productRepo.On("FindByID", isContext, args.id).Return(stub.findByID, stub.findByIDErr).Once()

			uc := usecase.NewProduct(productRepo)
			ucPdt, err := uc.View(context.Background(), args.id)
			assert.Equal(tc.wantErr, unwrapErr(err))
			assert.Equal(tc.want, ucPdt)
		})
	}
}

func TestProductUsecaseDelete(t *testing.T) {
	type stub struct {
		findByID    *domain.Product
		findByIDErr error
		deleteErr   error
	}

	type args struct {
		id     uuid.UUID
		userID uuid.UUID
	}

	p := factories.NewProduct()

	tests := []struct {
		name    string
		args    args
		stub    stub
		wantErr error
	}{
		{
			name:    "success",
			args:    args{id: p.ID, userID: p.UserID},
			stub:    stub{findByID: p},
			wantErr: nil,
		},
		{
			name:    "unauthorized",
			args:    args{id: p.ID, userID: uuid.New()},
			stub:    stub{findByID: p},
			wantErr: usecase.ErrProductUnauthorized,
		},
		{
			name:    "productRepo.findByID failed",
			args:    args{id: p.ID},
			stub:    stub{findByIDErr: errors.New("db error")},
			wantErr: errors.New("db error"),
		},
		{
			name:    "productRepo.delete failed",
			args:    args{id: p.ID, userID: p.UserID},
			stub:    stub{findByID: p, deleteErr: errors.New("db error")},
			wantErr: errors.New("db error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			args, stub := tc.args, tc.stub

			productRepo := new(mocks.ProductRepository)
			productRepo.On("FindByID", isContext, args.id).Return(stub.findByID, stub.findByIDErr).Once()
			productRepo.On("Delete", isContext, args.id).Return(stub.deleteErr).Once()

			uc := usecase.NewProduct(productRepo)
			err := uc.Delete(context.Background(), args.id, args.userID)
			assert.Equal(tc.wantErr, unwrapErr(err))
		})
	}
}

func TestProductUsecaseCreate(t *testing.T) {
	type stub struct {
		create    *domain.Product
		createErr error
	}

	type args usecase.CreateProductDto
	p := factories.NewProduct()

	tests := []struct {
		name    string
		args    args
		stub    stub
		want    *domain.Product
		wantErr error
	}{
		{
			name:    "success",
			args:    args{Name: "colorful socks", UserID: p.UserID},
			stub:    stub{create: p},
			want:    p,
			wantErr: nil,
		},
		{
			name:    "name bad input",
			args:    args{Name: "!@#$!@#", UserID: p.UserID},
			wantErr: usecase.ErrProductNameBadFormat,
		},
		{
			name:    "productRepo.create failed",
			args:    args{Name: "colorful socks", UserID: p.UserID},
			stub:    stub{createErr: errors.New("db error")},
			wantErr: errors.New("db error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			args, stub := tc.args, tc.stub

			productRepo := new(mocks.ProductRepository)
			productRepo.On("Create", isContext, args.Name, args.UserID).Return(stub.create, stub.createErr).Once()

			uc := usecase.NewProduct(productRepo)
			pdt, err := uc.Create(context.Background(), usecase.CreateProductDto(args))
			assert.Equal(tc.wantErr, unwrapErr(err))
			assert.Equal(tc.want, pdt)
		})
	}
}
