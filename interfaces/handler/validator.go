package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	validator *validator.Validate
}

func newValidator() (echo.Validator, error) {
	v := validator.New()
	if err := v.RegisterValidation("is-uuid", isValidUUID); err != nil {
		return nil, err
	}

	return &Validator{v}, nil
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func isValidUUID(fl validator.FieldLevel) bool {
	id, ok := fl.Field().Interface().(uuid.UUID)
	return ok && id != uuid.Nil
}
