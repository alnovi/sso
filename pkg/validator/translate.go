package validator

import (
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ruTranslations "github.com/go-playground/validator/v10/translations/ru"
)

func newTranslator(validate *validator.Validate) ut.Translator {
	ruTranslator := ru.New()
	universalTranslator := ut.New(ruTranslator, ruTranslator)
	translator, _ := universalTranslator.GetTranslator("ru")
	_ = ruTranslations.RegisterDefaultTranslations(validate, translator)

	return translator
}

func translateOverride(validation *validator.Validate, trans ut.Translator) {
	_ = validation.RegisterTranslation("datetime", trans, func(ut ut.Translator) error {
		return ut.Add("datetime", "{0} не верный формат", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("datetime", fe.Field())
		return t
	})
}
