package common

import (
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
)

func AbsolutePath(fl validator.FieldLevel) bool {
	return filepath.IsAbs(fl.Field().String())
}

func ISODate(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}

func ValidScenario(fl validator.FieldLevel) bool {
	for _, v := range []string{"NETWORK", "PROXY", "HTTP_RESPONSE"} {
		if v == fl.Field().String() {
			return true
		}
	}
	return false
}

func InitialiseValidator(validate *validator.Validate) *validator.Validate {
	validate.RegisterValidation("validScenario", ValidScenario, false)
	validate.RegisterValidation("absolutePath", AbsolutePath, false)
	validate.RegisterValidation("isodate", ISODate, false)
	return validate
}
