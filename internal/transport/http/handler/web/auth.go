package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/transport/http/handler"
	"github.com/alnovi/sso/internal/transport/http/request"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

type AuthUseCase interface {
	ValidateGrantType(ctx context.Context, inp dto.InpValidateGrantType) (*entity.Client, error)
	CodeByCredentials(ctx context.Context, inp dto.InpAuthByCredentials) (*entity.Token, error)
}

type AuthHandler struct {
	handler.BaseHandler
	clientId string
	uc       AuthUseCase
}

func NewAuthHandler(clientId string, uc AuthUseCase) *AuthHandler {
	return &AuthHandler{clientId: clientId, uc: uc}
}

func (h *AuthHandler) Home(c echo.Context) error {
	authURI := fmt.Sprintf("/oauth/authorize?response_type=code&client_id=%s", h.clientId)
	return c.Redirect(http.StatusFound, authURI)
}

func (h *AuthHandler) SignIn(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := h.uc.ValidateGrantType(ctx, dto.InpValidateGrantType{
		ClientID:    c.QueryParam("client_id"),
		GrantType:   c.QueryParam("response_type"),
		RedirectURI: c.QueryParam("redirect_uri"),
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	return c.Render(http.StatusOK, "auth.html", echo.Map{
		"Query":    c.Request().URL.RawQuery,
		"AppId":    client.ID,
		"AppName":  client.Name,
		"AppIcon":  client.Icon,
		"AppColor": client.Color,
		"AppImage": client.Image,
	})
}

func (h *AuthHandler) Authorize(c echo.Context) error {
	var err error
	var req = request.Authorize{}

	ctx := c.Request().Context()

	if err = c.Bind(&req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return err
	}

	client, err := h.uc.ValidateGrantType(ctx, dto.InpValidateGrantType{
		ClientID:    c.QueryParam("client_id"),
		GrantType:   c.QueryParam("response_type"),
		RedirectURI: c.QueryParam("redirect_uri"),
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	token, err := h.uc.CodeByCredentials(ctx, dto.InpAuthByCredentials{
		Client:   client,
		Email:    req.Login,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, exception.ErrUserNotFound) {
			return validator.NewValidateErrorWithMessage("login", "Пользователь не найден")
		}
		if errors.Is(err, exception.ErrPasswordIncorrect) {
			return validator.NewValidateErrorWithMessage("password", "Не верный пароль")
		}
		return err
	}

	callbackURI := fmt.Sprintf("%s?code=%s&state=%s", client.Callback, token.Hash, c.QueryParam("state"))

	if h.RequestIsJSON(c.Request()) {
		return c.JSON(http.StatusOK, echo.Map{"location": callbackURI})
	} else {
		return c.Redirect(http.StatusFound, callbackURI)
	}
}

func (h *AuthHandler) Token(c echo.Context) error {
	// TODO: implement me
	return nil
}

func (h *AuthHandler) Route(e *echo.Group) {
	e.GET("/", h.Home)
	e.GET("/oauth/authorize/", h.SignIn)
	e.POST("/oauth/authorize/", h.Authorize)
	e.POST("/oauth/token/", h.Token)
}
