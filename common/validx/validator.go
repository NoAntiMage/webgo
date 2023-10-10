package validx

import v "github.com/go-playground/validator/v10"

var Valid *v.Validate

func ValidatorInit() {
	Valid = v.New()
}

func GetValidator() *v.Validate {
	return Valid
}
