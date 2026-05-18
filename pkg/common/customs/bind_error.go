package customs

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
)

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "alpha":
		return "Only alphabet characters are allowed"
	case "gte":
		return "Must be greater than or equal to " + fe.Param()
	case "lte":
		return "Must be less than or equal to " + fe.Param()
	}

	return "Invalid value"
}

func HandleBindError(err error) ErrorValues {

	var errValues ErrorValues

	// type mismatch error
	var unmarshalTypeError *json.UnmarshalTypeError

	if errorsAs := errors.As(err, &unmarshalTypeError); errorsAs {
		errValues = append(errValues, ErrorValue{
			Key:   unmarshalTypeError.Field,
			Value: "Invalid type, expected " + unmarshalTypeError.Type.String(),
		})
		return errValues
	}

	// validation errors
	var validationErrors validator.ValidationErrors

	if errorsAs := errors.As(err, &validationErrors); errorsAs {
		for _, fe := range validationErrors {
			errValues = append(errValues, ErrorValue{
				Key:   fe.Field(),
				Value: getErrorMsg(fe),
			})
		}
		return errValues
	}

	// invalid JSON
	var syntaxError *json.SyntaxError

	if errorsAs := errors.As(err, &syntaxError); errorsAs {
		errValues = append(errValues, ErrorValue{
			Key:   "body",
			Value: "Invalid JSON format",
		})
		return errValues
	}

	// fallback
	errValues = append(errValues, ErrorValue{
		Key:   "body",
		Value: "Invalid request body",
	})

	return errValues
}
