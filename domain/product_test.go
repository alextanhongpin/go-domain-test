package domain_test

import (
	"testing"
	"time"

	"github.com/alextanhongpin/go-domain-test/domain"
	"github.com/alextanhongpin/go-domain-test/domain/factories"
	"github.com/alextanhongpin/go-domain-test/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProductIsPublished(t *testing.T) {
	tests := []struct {
		name string
		args *domain.Product
		want bool
	}{
		{
			name: "no published at",
			args: &domain.Product{},
			want: false,
		},
		{
			name: "published in the past",
			args: func() *domain.Product {
				pdt := &domain.Product{}
				pdt.PublishedAt = types.Ptr(time.Now().Add(-1 * time.Second))
				return pdt
			}(),
			want: true,
		},
		{
			name: "published in the future",
			args: func() *domain.Product {
				pdt := &domain.Product{}
				pdt.PublishedAt = types.Ptr(time.Now().Add(1 * time.Second))
				return pdt
			}(),
			want: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.args.IsPublished())
		})
	}
}

func TestProductIsMine(t *testing.T) {
	type args struct {
		userID uuid.UUID
	}

	p := factories.NewProduct()
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is mine",
			args: args{userID: p.UserID},
			want: true,
		},
		{
			name: "is not mine",
			args: args{userID: uuid.New()},
			want: false,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, p.IsMine(tc.args.userID))
		})
	}
}

func TestProductName(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{name: "valid", args: "colorful socks", want: true},
		{name: "invalid", args: "%!@", want: false},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			name := domain.ProductName(tc.args)
			assert.Equal(t, tc.want, name.Valid())
		})
	}
}
