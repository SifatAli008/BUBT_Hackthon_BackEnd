package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom tag name function to use JSON tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct validates a struct using the validator package
func ValidateStruct(s interface{}) []string {
	var errors []string

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, formatError(err))
		}
	}

	return errors
}

// formatError formats validation error messages
func formatError(err validator.FieldError) string {
	var message string

	field := err.Field()
	if err.Tag() == "required" {
		message = fmt.Sprintf("%s is required", field)
	} else {
		switch err.Tag() {
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters", field, err.Param())
		case "max":
			message = fmt.Sprintf("%s must be at most %s characters", field, err.Param())
		case "len":
			message = fmt.Sprintf("%s must be exactly %s characters", field, err.Param())
		case "oneof":
			message = fmt.Sprintf("%s must be one of: %s", field, err.Param())
		case "uuid":
			message = fmt.Sprintf("%s must be a valid UUID", field)
		case "url":
			message = fmt.Sprintf("%s must be a valid URL", field)
		case "numeric":
			message = fmt.Sprintf("%s must be numeric", field)
		case "alpha":
			message = fmt.Sprintf("%s must contain only letters", field)
		case "alphanum":
			message = fmt.Sprintf("%s must contain only letters and numbers", field)
		default:
			message = fmt.Sprintf("%s is invalid", field)
		}
	}

	return message
}

// ValidateEmail validates an email address
func ValidateEmail(email string) bool {
	return validate.Var(email, "required,email") == nil
}

// ValidateUUID validates a UUID string
func ValidateUUID(uuid string) bool {
	return validate.Var(uuid, "required,uuid") == nil
}
