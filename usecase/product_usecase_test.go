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
	matchCtx = mock.MatchedBy(func(ctx context.Context) bool {
		return true
	})
)

func TestProductUsecaseView(t *testing.T) {
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

	t.Run("success", func(t *testing.T) {
		f := newViewProductFlow()
		err := f.exec()
		assert.Nil(t, err)
	})

	t.Run("not yet published", func(t *testing.T) {
		f := newViewProductFlow()
		f.stub.findByID.data = productWithoutPublishedAt()
		assert.ErrorIs(t, f.exec(), usecase.ErrProductNotFound)
	})

	t.Run("published in the future", func(t *testing.T) {
		f := newViewProductFlow()
		f.stub.findByID.data = productPublishedInTheFuture()
		assert.ErrorIs(t, f.exec(), usecase.ErrProductNotFound)
	})

	t.Run("error when find by id", func(t *testing.T) {
		f := newViewProductFlow()
		f.stub.findByID.err = wantErr
		assert.ErrorIs(t, f.exec(), wantErr)
	})
}

func TestProductUsecaseDeleteFlow(t *testing.T) {
	wantErr := errors.New("want error")

	t.Run("success", func(t *testing.T) {
		f := newDeleteProductFlow()
		assert.Nil(t, f.exec())
	})

	t.Run("unauthorized user id", func(t *testing.T) {
		f := newDeleteProductFlow()
		f.args.userID = uuid.New()
		assert.ErrorIs(t, f.exec(), usecase.ErrProductUnauthorized)
	})

	t.Run("error when finding product by id", func(t *testing.T) {
		f := newDeleteProductFlow()
		f.stub.findByID.err = wantErr
		assert.ErrorIs(t, f.exec(), wantErr)
	})

	t.Run("error when delete", func(t *testing.T) {
		f := newDeleteProductFlow()
		f.stub.deleteErr = wantErr
		assert.ErrorIs(t, f.exec(), wantErr)
	})
}

func TestProductUsecaseDelete(t *testing.T) {
	// This test is just a comparison to the method above.
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
			name: "repo.findByID failed",
			stubFn: func(s *stub) {
				s.findByIDErr = wantErr
			},
			wantErr: wantErr,
		},
		{
			name: "repo.delete failed",
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

			repo := new(mocks.ProductRepository)
			repo.On("FindByID", matchCtx, args.id).Return(stub.findByID, stub.findByIDErr).Once()
			repo.On("Delete", matchCtx, args.id).Return(stub.deleteErr).Once()

			uc := usecase.NewProduct(repo)
			err := uc.Delete(context.Background(), args.id, args.userID)
			assert.ErrorIs(err, tc.wantErr)
			t.Logf("%s: %s\n", tc.name, err)
		})
	}
}

func TestProductUsecaseCreate(t *testing.T) {
	wantErr := errors.New("want error")

	t.Run("success", func(t *testing.T) {
		f := newCreateProductFlow()
		assert.Nil(t, f.exec())
	})

	t.Run("when input invalid name", func(t *testing.T) {
		f := newCreateProductFlow()
		f.args.Name = "!@#$!@#"

		assert.ErrorIs(t, f.exec(), usecase.ErrProductNameBadFormat)
	})

	t.Run("error when create", func(t *testing.T) {
		f := newCreateProductFlow()
		f.stub.create.err = wantErr
		assert.ErrorIs(t, f.exec(), wantErr)
	})
}

type result[T any] struct {
	data T
	err  error
}

type viewProductFlow struct {
	args struct {
		id uuid.UUID
	}
	stub struct {
		findByID result[*domain.Product]
	}
}

func newViewProductFlow() *viewProductFlow {
	p := factories.NewProduct()
	f := new(viewProductFlow)
	f.args.id = p.ID
	f.stub.findByID.data = p
	return f
}

func (f *viewProductFlow) exec() error {
	args := f.args
	stub := f.stub

	repo := new(mocks.ProductRepository)
	repo.On("FindByID", matchCtx, args.id).Return(stub.findByID.data, stub.findByID.err).Once()

	uc := usecase.NewProduct(repo)
	ctx := context.Background()
	_, err := uc.View(ctx, args.id)
	return err
}

type deleteProductFlow struct {
	args struct {
		id     uuid.UUID
		userID uuid.UUID
	}
	stub struct {
		findByID  result[*domain.Product]
		deleteErr error
	}
}

func newDeleteProductFlow() *deleteProductFlow {
	p := factories.NewProduct()

	f := new(deleteProductFlow)

	args := f.args
	args.id = p.ID
	args.userID = p.UserID

	stub := f.stub
	stub.findByID.data = p

	return f
}

func (f *deleteProductFlow) exec() error {
	ctx := context.Background()

	args := f.args
	stub := f.stub

	repo := new(mocks.ProductRepository)
	repo.On("FindByID", matchCtx, args.id).Return(stub.findByID.data, stub.findByID.err).Once()
	repo.On("Delete", matchCtx, args.id).Return(stub.deleteErr).Once()

	uc := usecase.NewProduct(repo)
	return uc.Delete(ctx, args.id, args.userID)
}

type createProductFlow struct {
	args usecase.CreateProductDto
	stub struct {
		create result[*domain.Product]
	}
}

func newCreateProductFlow() *createProductFlow {
	f := new(createProductFlow)

	f.args = usecase.CreateProductDto{
		Name:   "colorful socks",
		UserID: uuid.New(),
	}

	f.stub.create.data = factories.NewProduct()

	return f
}

func (f *createProductFlow) exec() error {
	args := f.args
	stub := f.stub

	repo := new(mocks.ProductRepository)
	repo.On("Create", matchCtx, args.Name, args.UserID).Return(stub.create.data, stub.create.err)

	ctx := context.Background()
	uc := usecase.NewProduct(repo)
	_, err := uc.Create(ctx, args)
	return err
}
