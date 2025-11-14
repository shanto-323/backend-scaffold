package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/backend-scaffold/internal/server/errs"
)

type Validatable interface {
	Validate() error
}

type CustomValidationErrors []CustomValidationError

type CustomValidationError struct {
	Field   string
	message string
}

func (c CustomValidationErrors) Error() string {
	return "Validation failed"
}

func BindAndValidate(ctx echo.Context, payload Validatable) error {
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	if msg, err := validateStruct(payload); err != nil {
		return errs.NewBadRequestError(msg, false, nil, err, nil)
	}

	return nil
}

func validateStruct(v Validatable) (string, []errs.FieldError) {
	if err := v.Validate(); err != nil {
		return extrectValidationErrors(err)
	}
	return "", nil
}

func extrectValidationErrors(err error) (string, []errs.FieldError) {
	var fieldErrors []errs.FieldError

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		customValidationError := err.(CustomValidationErrors)
		for _, err := range customValidationError {
			fieldErrors = append(fieldErrors, errs.FieldError{
				Field: err.Field,
				Error: err.message,
			})
		}
	}

	for _, err := range validationErrors {
		field := strings.ToLower(err.Field())
		var msg string

		switch err.Tag() {
		case "required":
			msg = "is required"
		case "min":
			if err.Type().Kind() == reflect.String {
				msg = fmt.Sprintf("must be at least %s characters", err.Param())
			} else {
				msg = fmt.Sprintf("must be at least %s", err.Param())
			}
		case "max":
			if err.Type().Kind() == reflect.String {
				msg = fmt.Sprintf("must not exceed %s characters", err.Param())
			} else {
				msg = fmt.Sprintf("must not exceed %s", err.Param())
			}
		case "oneof":
			msg = fmt.Sprintf("must be one of: %s", err.Param())
		case "email":
			msg = "must be a valid email address"
		case "e164":
			msg = "must be a valid phone number with country code"
		case "uuid":
			msg = "must be a valid UUID"
		case "uuidList":
			msg = "must be a comma-separated list of valid UUIDs"
		case "dive":
			msg = "some items are invalid"
		default:
			if err.Param() != "" {
				msg = fmt.Sprintf("%s: %s:%s", field, err.Tag(), err.Param())
			} else {
				msg = fmt.Sprintf("%s: %s", field, err.Tag())
			}
		}

		fieldErrors = append(fieldErrors, errs.FieldError{
			Field: strings.ToLower(err.Field()),
			Error: msg,
		})
	}

	return "Validation failed", fieldErrors
}
