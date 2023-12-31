package factories

import (
	"log"
	"time"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/alextanhongpin/go-domain-test/types"
	"github.com/google/uuid"
)

func NewProduct(scenarios ...string) *domain.Product {
	p := newProduct()

	for _, s := range scenarios {
		switch s {
		case "published_in_the_future":
			p.PublishedAt = types.Ptr(time.Now().Add(1 * time.Second))
		case "no_published_at":
			p.PublishedAt = nil
		default:
			log.Fatalf("unknown product scenario: %s", s)
		}
	}

	return p
}

func newProduct() *domain.Product {
	return &domain.Product{
		ID:          uuid.New(),
		Name:        "colorful socks",
		PublishedAt: types.Ptr(time.Now()),
		UserID:      uuid.New(),
	}
}

func NewUser() *domain.User {
	return &domain.User{
		ID:   uuid.New(),
		Name: "John Appleseed",
	}
}
