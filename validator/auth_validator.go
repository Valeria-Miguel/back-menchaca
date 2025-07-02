package validators

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateLoginInput(input interface{}) error {
	return validate.Struct(input)
}

func ValidateRefreshInput(input interface{}) error {
	return validate.Struct(input)
}