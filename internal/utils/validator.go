package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/jashakimov/converter/internal/api"
)

func NewValidator() *validator.Validate {
	val := validator.New()
	val.RegisterValidation("allowed-preset", func(fl validator.FieldLevel) bool {
		for i := range api.AllowedPresets {
			if api.AllowedPresets[i] == api.Preset(fl.Field().String()) {
				return true
			}
		}
		return false
	})

	return val
}
