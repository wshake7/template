package validator

import (
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func init() {
	validate.SetTagName("binding")
}

func Struct[T any](t T) error {
	return validate.Struct(t)
}

type StructValidator struct{}

func (v *StructValidator) Validate(out any) error {
	return validate.Struct(out)
}

type Validator interface {
	Validate() error
}
