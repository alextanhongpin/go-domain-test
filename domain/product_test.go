package domain_test

import (
	"sort"
	"testing"
	"time"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/alextanhongpin/go-domain-test/domain/factories"
	"github.com/alextanhongpin/go-domain-test/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProductIsPublished(t *testing.T) {
	newProduct := func(duration time.Duration) *domain.Product {
		pdt := &domain.Product{}
		pdt.PublishedAt = types.Ptr(time.Now().Add(duration))
		return pdt
	}

	tc := make(testcase)
	tc["no published at"] = (&domain.Product{}).IsPublished() == false
	tc["published in the past"] = newProduct(-1*time.Second).IsPublished() == true
	tc["published in the future"] = newProduct(1*time.Second).IsPublished() == false
	tc.Run(t)
}

func TestProductIsMine(t *testing.T) {
	p := factories.NewProduct()

	tc := make(testcase)
	tc["is mine"] = p.IsMine(p.UserID) == true
	tc["is not mine"] = p.IsMine(uuid.New()) == false
	tc.Run(t)
}

func TestProductName(t *testing.T) {
	tc := make(testcase)
	tc["valid"] = domain.ProductName("colorful stocks").Valid() == true
	tc["invalid"] = domain.ProductName("%!@").Valid() == false
	tc.Run(t)
}

type testcase map[string]bool

func (tc testcase) Run(t *testing.T) {
	keys := make([]string, 0, len(tc))
	for k := range tc {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		t.Run(k, func(t *testing.T) {
			assert.True(t, tc[k])
		})
	}
}
