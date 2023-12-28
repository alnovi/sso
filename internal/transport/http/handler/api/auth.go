package api

import (
	"errors"
	"net/http"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/pkg/cookies"
	"github.com/alnovi/sso/internal/transport/http/request"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/internal/usecase"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	auth   usecase.Auth
	client usecase.Client
}

func NewAuth(auth usecase.Auth, client usecase.Client) *Auth {
	return &Auth{
		auth:   auth,
		client: client,
	}
}

func (h *Auth) SignIn(c echo.Context) error {
	var err error
	var req request.SignInJson
	var ctx = c.Request().Context()

	if err = c.Bind(&req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return err
	}

	dtoClient := dto.ClientForAuth{
		ClientId:    c.QueryParam("client_id"),
		RedirectURI: c.QueryParam("redirect_uri"),
	}

	client, err := h.client.ClientForAuth(ctx, dtoClient)
	if err != nil {
		if exception.Is(err) {
			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	dtoAuth := dto.AuthByCredentials{
		Client:   *client,
		Login:    req.Login,
		Password: req.Password,
		IP:       c.RealIP(),
		Agent:    c.Request().UserAgent(),
	}

	user, callback, err := h.auth.AuthByCredentials(ctx, dtoAuth)
	if err != nil {
		if errors.Is(err, exception.UserNotFound) {
			return validator.NewValidateErrorWithMessage("login", "Пользователь не найден")
		}
		if errors.Is(err, exception.PasswordIncorrect) {
			return validator.NewValidateErrorWithMessage("password", "Не верный пароль")
		}
		return err
	}

	c.SetCookie(cookies.Auth(user.Id, req.IsRemember))

	return c.JSON(http.StatusOK, response.Location{Location: callback.String()})
}
