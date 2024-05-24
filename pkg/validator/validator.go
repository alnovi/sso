package validator

import (
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type EchoValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func NewEchoValidator() *EchoValidator {
	validate := validator.New()
	translator := newTranslator(validate)

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0] //nolint:gomnd
		if name == "" {
			name = field.Tag.Get("form")
		}

		if name == "-" {
			return ""
		}

		return name
	})

	return &EchoValidator{
		validator:  validate,
		translator: translator,
	}
}

func (e *EchoValidator) Validate(i interface{}) error {
	if err := e.validator.Struct(i); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		return e.convertErrors(err)
	}

	return nil
}

func (e *EchoValidator) AddValidTrans(tag, msg string, fn validator.Func) error {
	err := e.AddValidation(tag, fn)
	if err != nil {
		return err
	}

	err = e.AddTranslation(tag, msg)
	if err != nil {
		return err
	}

	return nil
}

func (e *EchoValidator) AddValidation(tag string, fn validator.Func) error {
	return e.validator.RegisterValidation(tag, fn, false)
}

func (e *EchoValidator) AddTranslation(tag, msg string) error {
	return e.validator.RegisterTranslation(tag, e.translator, func(ut ut.Translator) error {
		return ut.Add(tag, msg, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}

func (e *EchoValidator) convertErrors(validErrors error) *ValidateError {
	err := NewValidateError()

	for _, fe := range validErrors.(validator.ValidationErrors) {
		namespace := strings.SplitN(fe.Namespace(), ".", 2) //nolint:gomnd
		err.Fields[namespace[1]] = fe.Translate(e.translator)
	}

	return err
}
