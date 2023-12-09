package validator

import (
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	translator ut.Translator
	validate   *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New()
	translator := newTranslator(validate)
	translateOverride(validate, translator)

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		translator: translator,
		validate:   validate,
	}
}

func (v *Validator) Validate(s interface{}) error {
	if validErrors := v.validate.Struct(s); validErrors != nil {
		return v.convertErrors(validErrors)
	}

	return nil
}

func (v *Validator) convertErrors(validErrors error) *ValidateError {
	err := NewValidateError()

	for _, e := range validErrors.(validator.ValidationErrors) {
		namespace := strings.SplitN(e.Namespace(), ".", 2)
		err.Fields[namespace[1]] = e.Translate(v.translator)
	}

	return err
}
