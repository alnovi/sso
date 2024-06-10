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
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

type authUseCase interface {
	ValidateResponseType(ctx context.Context, inp dto.ValidateResponseType) (*entity.Client, error)
	CodeByCredentials(ctx context.Context, inp dto.AuthByCredentials) (*entity.Token, error)
	ClientByResetPassword(ctx context.Context, hash string) (*entity.Client, error)
	ForgotPassword(ctx context.Context, inp dto.ForgotPassword) (*entity.User, error)
	ResetPassword(ctx context.Context, inp dto.ResetPassword) (*entity.Client, error)
}

type AuthHandler struct {
	handler.BaseHandler
	clientId string
	uc       authUseCase
}

func NewAuthHandler(clientId string, uc authUseCase) *AuthHandler {
	return &AuthHandler{clientId: clientId, uc: uc}
}

func (h *AuthHandler) Home(c echo.Context) error {
	authURI := fmt.Sprintf("/oauth/authorize?response_type=%s&client_id=%s", dto.ResponseTypeCode, h.clientId)
	return c.Redirect(http.StatusFound, authURI)
}

func (h *AuthHandler) SignIn(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := h.uc.ValidateResponseType(ctx, dto.ValidateResponseType{
		ClientID:     c.QueryParam("client_id"),
		ResponseType: c.QueryParam("response_type"),
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

	client, err := h.uc.ValidateResponseType(ctx, dto.ValidateResponseType{
		ClientID:     c.QueryParam("client_id"),
		ResponseType: c.QueryParam("response_type"),
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	token, err := h.uc.CodeByCredentials(ctx, dto.AuthByCredentials{
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

func (h *AuthHandler) ForgotPasswordPage(c echo.Context) error {
	ctx := c.Request().Context()

	if c.Request().URL.RawQuery == "" {
		forgotPasswordUrl := fmt.Sprintf("/oauth/forgot-password?response_type=%s&client_id=%s", dto.ResponseTypeCode, h.clientId)
		return c.Redirect(http.StatusFound, forgotPasswordUrl)
	}

	client, err := h.uc.ValidateResponseType(ctx, dto.ValidateResponseType{
		ClientID:     c.QueryParam("client_id"),
		ResponseType: c.QueryParam("response_type"),
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

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var err error
	var req = request.ForgotPassword{}

	ctx := c.Request().Context()

	if err = c.Bind(&req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return err
	}

	client, err := h.uc.ValidateResponseType(ctx, dto.ValidateResponseType{
		ClientID:     c.QueryParam("client_id"),
		ResponseType: c.QueryParam("response_type"),
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	user, err := h.uc.ForgotPassword(ctx, dto.ForgotPassword{
		Client: client,
		Email:  req.Login,
		IP:     c.RealIP(),
		Agent:  c.Request().UserAgent(),
	})

	if err != nil {
		if errors.Is(err, exception.ErrUserNotFound) {
			return validator.NewValidateErrorWithMessage("login", "Пользователь не найден")
		}
		return err
	}

	return c.JSON(http.StatusOK, response.Message{
		Message: "Ссылка, для сброса пароля, отправлена на " + user.Email,
	})
}

func (h *AuthHandler) ResetPasswordPage(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := h.uc.ClientByResetPassword(ctx, c.QueryParam("hash"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	return c.Render(http.StatusOK, "auth.html", echo.Map{
		"Query":    fmt.Sprintf("/oauth/authorize?response_type=%s&client_id=%s", dto.ResponseTypeCode, client.ID),
		"AppId":    client.ID,
		"AppName":  client.Name,
		"AppIcon":  client.Icon,
		"AppColor": client.Color,
		"AppImage": client.Image,
	})
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var err error
	var req = request.ResetPassword{}

	ctx := c.Request().Context()

	if err = c.Bind(&req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return err
	}

	client, err := h.uc.ResetPassword(ctx, dto.ResetPassword{
		Hash:     req.Hash,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, exception.ErrTokenNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
		}
		if errors.Is(err, exception.ErrClientNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
		}
		if errors.Is(err, exception.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	authURI := fmt.Sprintf("/oauth/authorize?response_type=%s&client_id=%s", dto.ResponseTypeCode, client.ID)

	if h.RequestIsJSON(c.Request()) {
		return c.JSON(http.StatusOK, echo.Map{"location": authURI})
	} else {
		return c.Redirect(http.StatusFound, authURI)
	}
}

func (h *AuthHandler) Route(e *echo.Group) {
	e.GET("/", h.Home)
	e.GET("/oauth/authorize/", h.SignIn)
	e.POST("/oauth/authorize/", h.Authorize)
	e.GET("/oauth/forgot-password/", h.ForgotPasswordPage)
	e.POST("/oauth/forgot-password/", h.ForgotPassword)
	e.GET("/oauth/reset-password/", h.ResetPasswordPage)
	e.POST("/oauth/reset-password/", h.ResetPassword)
}
