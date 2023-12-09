package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/transport/http/request"
	"github.com/alnovi/sso/internal/usecase"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	auth   usecase.Auth
	client usecase.Client
}

func NewAuthHandler(auth usecase.Auth, client usecase.Client) *AuthHandler {
	return &AuthHandler{auth: auth, client: client}
}

// TODO: добавить документацию
func (h *AuthHandler) AuthForm(c echo.Context) error {
	var err error

	ctx := c.Request().Context()

	dtoClient := dto.ClientForAuth{
		ClientId:    c.QueryParam("client_id"),
		RedirectURI: c.QueryParam("redirect_uri"),
	}

	client, err := h.client.ClientForAuth(ctx, dtoClient)
	if exception.Is(err) {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	if cookie, err := c.Cookie("uid"); err == nil {
		dtoAuth := dto.AuthById{
			Client: *client,
			UserId: cookie.Value,
			IP:     c.RealIP(),
			Agent:  c.Request().UserAgent(),
		}

		_, callback, err := h.auth.AuthById(ctx, dtoAuth)
		if errors.Is(err, exception.AccessDenied) {
			return err
		}
		if err != nil {
			cookie.Expires = time.Now()
			c.SetCookie(cookie)
		}
		if err == nil {
			return c.Redirect(http.StatusFound, callback.String())
		}
	}

	return c.Render(http.StatusOK, "signin.html", echo.Map{
		"AppName":  client.Name,
		"AppLogo":  client.Logo,
		"AppImage": client.Image,
	})
}

// TODO: добавить документацию
func (h *AuthHandler) SignIn(c echo.Context) error {
	var err error
	var req request.SignIn
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
	if exception.Is(err) {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if err != nil {
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
	if errors.Is(err, exception.UserNotFound) {
		return validator.NewValidateErrorWithMessage("login", "Пользователь не найден")
	}
	if errors.Is(err, exception.PasswordIncorrect) {
		return validator.NewValidateErrorWithMessage("password", "Не верный пароль")
	}
	if err != nil {
		return err
	}

	if req.IsRemember {
		cookie := http.Cookie{
			Name:     "uid",
			Value:    user.Id,
			Path:     "/",
			Expires:  time.Now().AddDate(0, 1, 0),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		c.SetCookie(&cookie)
	}

	if c.Request().Header.Get("Content-Type") == "application/json" {
		return c.JSON(http.StatusOK, echo.Map{"location": callback.String()})
	}

	return c.Redirect(http.StatusFound, callback.String())
}
