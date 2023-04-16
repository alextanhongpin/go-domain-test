package domain

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

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
