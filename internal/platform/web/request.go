package web

import (
	"encoding/json"
	"errors"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"
	"reflect"
	"strings"
)

var validate = validator.New()

func init() {

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(val); err != nil {
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fields []FieldError

		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
			}
			fields = append(fields, field)
		}
		return &Error{
			Err:    errors.New("field validator error"),
			Status: http.StatusBadRequest,
			Fields: fields,
		}
	}

	return nil
}
