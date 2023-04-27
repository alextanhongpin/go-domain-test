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
)

func TestProductUsecaseView(t *testing.T) {
	type stub struct {
		findByID    *domain.Product
		findByIDErr error
	}

	productWithoutPublishedAt := func() *domain.Product {
		p := factories.NewProduct()
		p.PublishedAt = nil
		return p
	}

	productPublishedInTheFuture := func() *domain.Product {
		p := factories.NewProduct()
		p.PublishedAt = types.Ptr(time.Now().Add(1 * time.Second))
		return p
	}

	// This acts as sentinel error.
	wantErr := errors.New("want error")

	tests := []struct {
		name    string
		stubFn  func(*stub)
		wantErr error
	}{
		{
			name: "success",
		},
		{
			name: "not yet published",
			stubFn: func(s *stub) {
				s.findByID = productWithoutPublishedAt()
			},
			wantErr: usecase.ErrProductNotFound,
		},
		{
			name: "published in the future",
			stubFn: func(s *stub) {
				s.findByID = productPublishedInTheFuture()
			},
			wantErr: usecase.ErrProductNotFound,
		},
		{
			name: "failed",
			stubFn: func(s *stub) {
				s.findByIDErr = wantErr
			},
			wantErr: wantErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			p := factories.NewProduct()
			stub := stub{findByID: p}
			if tc.stubFn != nil {
				tc.stubFn(&stub)
			}

			productRepo := new(mocks.ProductRepository)
			productRepo.On("FindByID", isContext, p.ID).Return(stub.findByID, stub.findByIDErr).Once()

			uc := usecase.NewProduct(productRepo)
			ucPdt, err := uc.View(context.Background(), p.ID)
			if err != nil {
				assert.ErrorIs(err, tc.wantErr)
				assert.Nil(ucPdt)
				t.Logf("%s: %s\n", tc.name, err)
			} else {
				assert.Equal(p, ucPdt)
			}
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

	wantErr := errors.New("want error")

	tests := []struct {
		name    string
		argsFn  func(args) args
		stubFn  func(*stub)
		wantErr error
	}{
		{
			name:    "success",
			wantErr: nil,
		},
		{
			name: "unauthorized",
			argsFn: func(a args) args {
				a.userID = uuid.New()
				return a
			},
			wantErr: usecase.ErrProductUnauthorized,
		},
		{
			name: "productRepo.findByID failed",
			stubFn: func(s *stub) {
				s.findByIDErr = wantErr
			},
			wantErr: wantErr,
		},
		{
			name: "productRepo.delete failed",
			stubFn: func(s *stub) {
				s.deleteErr = wantErr
			},
			wantErr: wantErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			p := factories.NewProduct()
			args := args{
				id:     p.ID,
				userID: p.UserID,
			}
			if tc.argsFn != nil {
				args = tc.argsFn(args)
			}

			stub := stub{
				findByID: p,
			}
			if tc.stubFn != nil {
				tc.stubFn(&stub)
			}

			productRepo := new(mocks.ProductRepository)
			productRepo.On("FindByID", isContext, args.id).Return(stub.findByID, stub.findByIDErr).Once()
			productRepo.On("Delete", isContext, args.id).Return(stub.deleteErr).Once()

			uc := usecase.NewProduct(productRepo)
			err := uc.Delete(context.Background(), args.id, args.userID)
			assert.ErrorIs(err, tc.wantErr)
			t.Logf("%s: %s\n", tc.name, err)
		})
	}
}

func TestProductUsecaseCreate(t *testing.T) {
	type stub struct {
		create    *domain.Product
		createErr error
	}

	type args usecase.CreateProductDto

	wantErr := errors.New("want error")

	tests := []struct {
		name    string
		argsFn  func(args) args
		stubFn  func(s *stub)
		wantErr error
	}{
		{
			name:    "success",
			wantErr: nil,
		},
		{
			name: "name bad input",
			argsFn: func(a args) args {
				a.Name = "!@#$!@#"
				return a
			},
			wantErr: usecase.ErrProductNameBadFormat,
		},
		{
			name: "productRepo.create failed",
			stubFn: func(s *stub) {
				s.createErr = wantErr
			},
			wantErr: wantErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			p := factories.NewProduct()
			args := args{
				Name:   "colorful socks",
				UserID: p.UserID,
			}
			if tc.argsFn != nil {
				args = tc.argsFn(args)
			}

			stub := stub{create: p}
			if tc.stubFn != nil {
				tc.stubFn(&stub)
			}

			productRepo := new(mocks.ProductRepository)
			productRepo.On("Create", isContext, args.Name, args.UserID).Return(stub.create, stub.createErr).Once()

			uc := usecase.NewProduct(productRepo)
			pdt, err := uc.Create(context.Background(), usecase.CreateProductDto(args))
			if err != nil {
				assert.ErrorIs(err, tc.wantErr)
				assert.Nil(pdt)
				t.Logf("%s: %s\n", tc.name, err)
			} else {
				assert.Equal(pdt, p)
			}
		})
	}
}
