package usecase

import (
	_ "embed"

	"github.com/BurntSushi/toml"
	"github.com/alextanhongpin/go-core-microservice/types/errors"
)

var (
	//go:embed errors.toml
	errorBytes []byte
	_          = errors.MustLoad(errorBytes, toml.Unmarshal)

	// Product errors.
	ErrProductNotFound     = errors.Get("product.not_found")
	ErrProductUnauthorized = errors.Get("product.unauthorized")
)
