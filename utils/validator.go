package utils

import (
	"github.com/go-playground/validator/v10"
)

func Validate[T any](element *T) error {
	validator := validator.New()
	return validator.Struct(element)
}
