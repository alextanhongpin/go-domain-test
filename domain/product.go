package domain

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

var ErrNegativePrice = errors.New("-tive price")

var regexpProductName = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)

// regexp that matches alphanumeric only and spaces

type ProductName string

func (p ProductName) Valid() bool {
	return regexpProductName.MatchString(string(p))
}

type Product struct {
	ID          uuid.UUID
	Name        ProductName
	PublishedAt *time.Time
	UserID      uuid.UUID
	Price       int
}

func (p *Product) IsPublished() bool {
	if p.PublishedAt == nil {
		return false
	}

	return p.PublishedAt.Before(time.Now())
}

func (p *Product) IsMine(userID uuid.UUID) bool {
	return p.UserID == userID
}

func (p *Product) WithDiscount(discounts ...Discount) (*Product, error) {
	pc := *p

	for _, d := range discounts {
		pc.Price += d.Amount
	}

	if pc.Price < 0 {
		return p, ErrNegativePrice
	}

	return &pc, nil
}
