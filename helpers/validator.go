package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator(validator *validator.Validate) *Validator {
	validator.RegisterValidation("weakpassword", weakPassword)
	return &Validator{
		validator: validator,
	}
}

func (cv *Validator) Validate(i interface{}) error {
	// Validate the struct
	err := cv.validator.Struct(i)
	if err != nil {
		var errMsgs []string
		for _, e := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' is invalid: %s", e.Field(), e.ActualTag()))
		}
		return fmt.Errorf(strings.Join(errMsgs, "; "))
	}
	return nil
}

// Minimum 8 characters, contains at least one uppercase letter, one digit, and one special character
func weakPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",<>\./?\\|]`).MatchString(password)

	return hasUpper && hasDigit && hasSpecial
}
