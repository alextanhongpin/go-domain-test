package factories

import (
	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/google/uuid"
)

func NewUser(variants ...string) *domain.User {
	// Random user.
	u := &domain.User{
		ID:   uuid.New(),
		Name: "John Appleseed",
	}

	for _, v := range variants {
		switch v {
		case "john":
			u.ID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		}
	}

	return u
}
