// prefer fixtures?
package factories

import (
	"log"
	"time"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/alextanhongpin/go-domain-test/types"
	"github.com/google/uuid"
)

func NewProduct(variants ...string) *domain.Product {
	// Valid product.
	p := &domain.Product{
		ID:          uuid.New(),
		Name:        "colorful socks",
		PublishedAt: types.Ptr(time.Now()),
		UserID:      NewUser("john").ID, // Belongs to John.
		Price:       10,
	}

	for _, v := range variants {
		switch v {
		case "published_in_the_future":
			p.PublishedAt = types.Ptr(time.Now().Add(1 * time.Second))
		case "published":
			p.PublishedAt = types.Ptr(time.Now().Add(-1 * time.Second))
		case "no_published_at":
			p.PublishedAt = nil
		// TODO:
		case "chair":
		// implement specific product.
		// If there are nested entities, we can use the factory to create them, e.g. with_discount.
		case "unknown_user":
			p.UserID = uuid.New()
		default:
			log.Fatalf("unknown Product variant: %s", v)
		}
	}

	return p
}
