package usecase

import (
	"github.com/alextanhongpin/errcodes"
)

// Product errors.
var (
	ErrProductNotFound      = errcodes.New(errcodes.NotFound, "product_not_found", "Product does not exist or may have been deleted.")
	ErrProductUnauthorized  = errcodes.New(errcodes.Unauthorized, "product_unauthorized", "You do not have access to this product")
	ErrProductNameBadFormat = errcodes.New(errcodes.BadRequest, "product_name_bad_format", "Product name can only contain alphanumeric characters and spaces.")
)
