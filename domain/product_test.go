package domain_test

import (
	"testing"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/alextanhongpin/go-domain-test/domain/factories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProductIsPublished(t *testing.T) {
	as := assert.New(t)
	as.False(factories.NewProduct("no_published_at").IsPublished())
	as.True(factories.NewProduct("published").IsPublished())
	as.False(factories.NewProduct("published_in_the_future").IsPublished())
}

func TestProductIsMine(t *testing.T) {
	p := factories.NewProduct()

	as := assert.New(t)
	as.True(p.IsMine(p.UserID))
	as.False(p.IsMine(uuid.New()))
}

func TestProductName(t *testing.T) {
	as := assert.New(t)
	as.True(domain.ProductName("colorful stocks").Valid())
	as.False(domain.ProductName("%!@").Valid())
}

func TestProductDiscount(t *testing.T) {
	t.Run("+tive price after discount", func(t *testing.T) {
		d := factories.NewDiscount()
		p := factories.NewProduct()
		as := assert.New(t)
		p, err := p.WithDiscount(*d)
		as.Nil(err)
		as.Equal(5, p.Price)
	})

	t.Run("0 price after discount", func(t *testing.T) {
		d := factories.NewDiscount()
		p := factories.NewProduct()
		as := assert.New(t)
		p, err := p.WithDiscount(*d, *d)
		as.Nil(err)
		as.Equal(0, p.Price)
	})

	t.Run("-tive price after discount", func(t *testing.T) {
		d := factories.NewDiscount()
		p := factories.NewProduct()
		as := assert.New(t)
		p, err := p.WithDiscount(*d, *d, *d)
		as.ErrorIs(err, domain.ErrNegativePrice)
		as.Equal(10, p.Price)
	})
}
