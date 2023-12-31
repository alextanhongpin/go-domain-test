package domain_test

import (
	"testing"

	"github.com/alextanhongpin/go-domain-test/domain/factories"
	"github.com/stretchr/testify/assert"
)

func TestDiscount(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dis := factories.NewDiscount()
		assert.True(t, dis.IsValid())
	})

	t.Run("invalid", func(t *testing.T) {
		assert.False(t, factories.NewDiscount("positive_amount").IsValid())
		assert.False(t, factories.NewDiscount("negative_purchase_qty").IsValid())
	})
}
