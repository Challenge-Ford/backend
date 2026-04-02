package apperr

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Kind int

const (
	KindBadRequest Kind = iota
	KindValidation
	KindNotFound
	KindConflict
	KindForbidden
	KindUnauthorized
	KindInternal
)

func (k Kind) HTTPStatus() int {
	switch k {
	case KindBadRequest, KindValidation:
		return http.StatusBadRequest
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	case KindForbidden:
		return http.StatusForbidden
	case KindUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func (k Kind) String() string {
	switch k {
	case KindBadRequest:
		return "BAD_REQUEST"
	case KindValidation:
		return "VALIDATION_ERROR"
	case KindNotFound:
		return "NOT_FOUND"
	case KindConflict:
		return "CONFLICT"
	case KindForbidden:
		return "FORBIDDEN"
	case KindUnauthorized:
		return "UNAUTHORIZED"
	default:
		return "INTERNAL_ERROR"
	}
}

type Error struct {
	Kind             Kind
	Message          string
	Detail           any
	ValidationErrors []ValidationItem
	err              error
}

type ValidationItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.err)
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.err }

func (e *Error) WithDetail(detail any) *Error {
	e.Detail = detail
	return e
}

func BadRequest(msg string) *Error {
	return &Error{Kind: KindBadRequest, Message: msg}
}

func NotFound(resource string) *Error {
	return &Error{Kind: KindNotFound, Message: fmt.Sprintf("%s not found", resource)}
}

func Conflict(msg string) *Error {
	return &Error{Kind: KindConflict, Message: msg}
}

func Forbidden(msg string) *Error {
	return &Error{Kind: KindForbidden, Message: msg}
}

func Unauthorized(msg string) *Error {
	return &Error{Kind: KindUnauthorized, Message: msg}
}

func Internal(msg string, err error) *Error {
	return &Error{Kind: KindInternal, Message: msg, err: err}
}

func Validation(field, msg string) *Error {
	return &Error{
		Kind:    KindValidation,
		Message: msg,
		ValidationErrors: []ValidationItem{
			{Field: field, Message: msg},
		},
	}
}

func FromValidatorErr(err error) *Error {
	if err == nil {
		return nil
	}

	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) && len(verrs) > 0 {
		items := make([]ValidationItem, 0, len(verrs))
		for _, fe := range verrs {
			items = append(items, ValidationItem{
				Field:   fieldPath(fe),
				Message: humanizeMessage(fe),
			})
		}
		return &Error{
			Kind:             KindValidation,
			Message:          "validation failed",
			ValidationErrors: items,
		}
	}

	return &Error{Kind: KindValidation, Message: "invalid input"}
}

func humanizeMessage(fe validator.FieldError) string {
	tag := strings.ToLower(fe.Tag())
	param := fe.Param()

	if isNumeric(fe.Kind()) {
		switch tag {
		case "min", "gte":
			return fmt.Sprintf("must be greater than or equal to %s", param)
		case "max", "lte":
			return fmt.Sprintf("must be less than or equal to %s", param)
		}
	}

	switch tag {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must have at least %s characters", param)
	case "max":
		return fmt.Sprintf("must have at most %s characters", param)
	case "len":
		return fmt.Sprintf("must have exactly %s characters", param)
	case "vin":
		return "must be a valid 17-character VIN (I, O and Q are not allowed)"
	case "plate":
		return "must be a valid Brazilian plate (e.g. ABC-1234 or ABC1D23)"
	case "hexcolor":
		return "must be a valid hex color code (e.g. #FF0000)"
	default:
		return fmt.Sprintf("invalid value (rule '%s')", tag)
	}
}

func isNumeric(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func fieldPath(fe validator.FieldError) string {
	ns := fe.StructNamespace()
	if ns == "" {
		return toLowerCamel(fe.Field())
	}

	parts := strings.SplitN(ns, ".", 2)
	if len(parts) == 2 {
		ns = parts[1]
	}

	segments := strings.Split(ns, ".")
	for i, s := range segments {
		if !strings.Contains(s, "[") {
			segments[i] = toLowerCamel(s)
		}
	}
	return strings.Join(segments, ".")
}

func toLowerCamel(s string) string {
	if s == "" {
		return s
	}
	// All-uppercase words (e.g. VIN, ID) become fully lowercase
	if s == strings.ToUpper(s) {
		return strings.ToLower(s)
	}
	runes := []rune(s)
	runes[0] = []rune(strings.ToLower(string(runes[0])))[0]
	return string(runes)
}
