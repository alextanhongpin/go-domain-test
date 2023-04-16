package domain_test

import (
	"testing"
	"time"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/alextanhongpin/go-domain-test/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProduct(t *testing.T) {
	t.Run("no published at", func(t *testing.T) {
		pdt := &domain.Product{}
		assert.False(t, pdt.IsPublished())
	})

	t.Run("published in the past", func(t *testing.T) {
		pdt := &domain.Product{}
		pdt.PublishedAt = types.Ptr(time.Now().Add(-1 * time.Second))
		assert.True(t, pdt.IsPublished())
	})

	t.Run("published in the future", func(t *testing.T) {
		pdt := &domain.Product{}
		pdt.PublishedAt = types.Ptr(time.Now().Add(1 * time.Second))
		assert.False(t, pdt.IsPublished())
	})

	t.Run("is mine", func(t *testing.T) {
		userID := uuid.New()
		pdt := &domain.Product{UserID: userID}
		assert.True(t, pdt.IsMine(userID))
	})

	t.Run("is not mine", func(t *testing.T) {
		pdt := &domain.Product{UserID: uuid.New()}
		assert.False(t, pdt.IsMine(uuid.New()))
	})
}
