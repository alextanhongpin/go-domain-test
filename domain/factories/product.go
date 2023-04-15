package factories

import (
	"time"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/alextanhongpin/go-domain-test/types"
	"github.com/google/uuid"
)

func NewProduct() *domain.Product {
	return &domain.Product{
		ID:          uuid.New(),
		Name:        "colorful socks",
		PublishedAt: types.Ptr(time.Now()),
	}
}
