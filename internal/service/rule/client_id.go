package rule

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ClientID - use validate:"client_id"
type ClientID struct {
	regex *regexp.Regexp
}

func NewClientID() *ClientID {
	return &ClientID{regex: regexp.MustCompile(`^[a-z0-9-]+$`)}
}

func (r *ClientID) Tag() string {
	return "client_id"
}

func (r *ClientID) ErrMsg() string {
	return "Значение может содержать только буквы (в нижнем регистре), цифры и дефис"
}

func (r *ClientID) CallIfNull() bool {
	return true
}

func (r *ClientID) Validate(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}
	return r.regex.MatchString(fl.Field().String())
}
