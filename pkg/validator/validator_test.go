package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var (
	v = NewEchoValidator()
)

func TestValidator_Validate_1(t *testing.T) {
	type Data struct {
		Id    string `json:"-" validate:"required"`
		Email string `json:"email" validate:"required,min=3,email"`
	}

	data := Data{Id: "", Email: "example.com"}

	actual := v.Validate(data)
	expected := &ValidateError{Fields: map[string]string{
		"Id":    "Id обязательное поле",
		"email": "email должен быть email адресом",
	}}

	assert.Equal(t, expected, actual)
}

func TestValidator_Validate_2(t *testing.T) {
	data := struct {
		Id    string `json:"" validate:"required"`
		Email string `json:"email" validate:"required,min=3,email"`
	}{
		Id: "", Email: "example.com",
	}

	actual := v.Validate(data)
	expected := &ValidateError{Fields: map[string]string{
		"Id":    "Id обязательное поле",
		"email": "email должен быть email адресом",
	}}

	assert.Equal(t, expected, actual)
}

func TestValidator_AddValidationTranslation(t *testing.T) {
	data := struct {
		Tag string `json:"tag" validate:"tag"`
	}{
		Tag: "system",
	}

	err := v.AddValidationTranslation("tag", "не верный формат для тега", func(fl validator.FieldLevel) bool {
		return false
	})

	assert.NoError(t, err)

	expected := &ValidateError{Fields: map[string]string{"tag": "не верный формат для тега"}}
	actual := v.Validate(data)

	assert.Equal(t, expected, actual)
}
