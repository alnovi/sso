package handler

import (
	"errors"
	"net/http"

	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/pkg/utils"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

type Error struct {
}

func NewError() *Error {
	return &Error{}
}

func (h *Error) ErrorHandle(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	data := response.Error{
		Code:    http.StatusInternalServerError,
		Message: h.httpStatusText(http.StatusInternalServerError),
	}

	var echoHttpError *echo.HTTPError
	if errors.As(err, &echoHttpError) {
		data.Code = echoHttpError.Code
		data.Message = echoHttpError.Message.(string)

		if data.Message == http.StatusText(data.Code) {
			data.Message = h.httpStatusText(data.Code)
		}
	}

	var validateError *validator.ValidateError
	if errors.As(err, &validateError) {
		data.Code = http.StatusUnprocessableEntity
		data.Message = h.httpStatusText(http.StatusUnprocessableEntity)
		data.Validate = validateError.Fields
	}

	if errors.Is(err, exception.ErrClientAccessDenied) {
		data.Code = http.StatusForbidden
		data.Message = h.httpStatusText(http.StatusForbidden)
	}

	if errors.Is(err, exception.ErrNotAuthorization) {
		_ = c.Redirect(http.StatusFound, "/oauth/signin")
		return
	}

	if utils.RequestIsJson(c.Request()) {
		_ = c.JSON(data.Code, data)
	} else {
		_ = c.Render(data.Code, "error.html", data)
	}
}

func (h *Error) httpStatusText(code int) string {
	switch code {
	case http.StatusContinue:
		return "Продолжайте"
	case http.StatusSwitchingProtocols:
		return "Переключение протоколов"
	case http.StatusProcessing:
		return "Идёт обработка"
	case http.StatusEarlyHints:
		return "Ранняя метаинформация"
	case http.StatusOK:
		return "Успешно"
	case http.StatusCreated:
		return "Создано"
	case http.StatusAccepted:
		return "Принято"
	case http.StatusNonAuthoritativeInfo:
		return "Информация не авторитетна"
	case http.StatusNoContent:
		return "Нет содержимого"
	case http.StatusResetContent:
		return "Сбросить содержимое"
	case http.StatusPartialContent:
		return "Частичное содержимое"
	case http.StatusMultiStatus:
		return "Многостатусный"
	case http.StatusAlreadyReported:
		return "Уже сообщалось"
	case http.StatusIMUsed:
		return "Использовано IM"
	case http.StatusMultipleChoices:
		return "Множество выборов"
	case http.StatusMovedPermanently:
		return "Перемещено навсегда"
	case http.StatusFound:
		return "Перемещено временно"
	case http.StatusSeeOther:
		return "Смотреть другое"
	case http.StatusNotModified:
		return "Не изменялось"
	case http.StatusUseProxy:
		return "Использовать прокси"
	case http.StatusTemporaryRedirect:
		return "Временное перенаправление"
	case http.StatusPermanentRedirect:
		return "Постоянное перенаправление"
	case http.StatusBadRequest:
		return "Некорректный запрос"
	case http.StatusUnauthorized:
		return "Не авторизован"
	case http.StatusPaymentRequired:
		return "Необходима оплата"
	case http.StatusForbidden:
		return "Доступ запрещен"
	case http.StatusNotFound:
		return "Не найдено"
	case http.StatusMethodNotAllowed:
		return "Метод не поддерживается"
	case http.StatusNotAcceptable:
		return "Неприемлемо"
	case http.StatusProxyAuthRequired:
		return "Необходима аутентификация прокси"
	case http.StatusRequestTimeout:
		return "Истекло время ожидания"
	case http.StatusConflict:
		return "Конфликт"
	case http.StatusGone:
		return "Удалён"
	case http.StatusLengthRequired:
		return "Необходима длина"
	case http.StatusPreconditionFailed:
		return "Условие ложно"
	case http.StatusRequestEntityTooLarge:
		return "Полезная нагрузка слишком велика"
	case http.StatusRequestURITooLong:
		return "URI слишком длинный"
	case http.StatusUnsupportedMediaType:
		return "Неподдерживаемый тип данных"
	case http.StatusRequestedRangeNotSatisfiable:
		return "Диапазон не достижим"
	case http.StatusExpectationFailed:
		return "Ожидание не удалось"
	case http.StatusTeapot:
		return "Я чайник"
	case http.StatusMisdirectedRequest:
		return "Неверно направленный запрос"
	case http.StatusUnprocessableEntity:
		return "Ошибка ввода данных"
	case http.StatusLocked:
		return "Заблокировано"
	case http.StatusFailedDependency:
		return "Невыполненная зависимость"
	case http.StatusTooEarly:
		return "Слишком рано"
	case http.StatusUpgradeRequired:
		return "Необходимо обновление"
	case http.StatusPreconditionRequired:
		return "Необходимо предусловие"
	case http.StatusTooManyRequests:
		return "Слишком много запросов"
	case http.StatusRequestHeaderFieldsTooLarge:
		return "Поля заголовка запроса слишком большие"
	case http.StatusUnavailableForLegalReasons:
		return "Недоступно по юридическим причинам"
	case http.StatusInternalServerError:
		return "Ошибка сервера"
	case http.StatusNotImplemented:
		return "Не реализовано"
	case http.StatusBadGateway:
		return "Ошибочный шлюз"
	case http.StatusServiceUnavailable:
		return "Сервис недоступен"
	case http.StatusGatewayTimeout:
		return "Шлюз не отвечает"
	case http.StatusHTTPVersionNotSupported:
		return "Версия HTTP не поддерживается"
	case http.StatusVariantAlsoNegotiates:
		return "Вариант тоже проводит согласование"
	case http.StatusInsufficientStorage:
		return "Переполнение хранилища"
	case http.StatusLoopDetected:
		return "Обнаружено бесконечное перенаправление"
	case http.StatusNotExtended:
		return "Не расширено"
	case http.StatusNetworkAuthenticationRequired:
		return "Требуется сетевая аутентификация"
	default:
		return ""
	}
}
