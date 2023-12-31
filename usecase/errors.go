package usecase

import (
	"github.com/alextanhongpin/errors/causes"
	"github.com/alextanhongpin/errors/codes"
)

var (
	// Product errors.
	ErrProductNotFound      = causes.New(codes.NotFound, "product_not_found", "Product does not exist or may have been deleted.")
	ErrProductUnauthorized  = causes.New(codes.Unauthorized, "product_unauthorized", "You do not have access to this product")
	ErrProductNameBadFormat = causes.New(codes.BadRequest, "product_name_bad_format", "Product name can only contain alphanumeric characters and spaces.")

	// Discount errors.
	ErrDiscountInvalid = causes.New(codes.PreconditionFailed, "discount_invalid", "The discount cannot be applied")
)
