package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	PublishedAt *time.Time
}

func (p *Product) IsPublished() bool {
	if p.PublishedAt == nil {
		return false
	}

	return p.PublishedAt.Before(time.Now())
}
