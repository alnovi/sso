package handler

import (
	"errors"
	"net/http"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/internal/usecase"
	"github.com/labstack/echo/v4"
)

type TokenHandler struct {
	client usecase.Client
	token  usecase.Token
}

func NewTokenHandler(client usecase.Client, token usecase.Token) *TokenHandler {
	return &TokenHandler{client: client, token: token}
}

// GenerateToken godoc
// @ID          GenerateToken
// @Summary     Получение токена доступа
// @Description Получение токена доступа по коду или обновление токена доступа
// @Tags        Авторизация
// @Produce     json
// @Param       client_id query string true "ID клиента"
// @Param       client_secret query string true "Секрет клиента"
// @Param       code query string false "Единоразовый авторизационный код, если grant_type=authorization_code"
// @Param       refresh_token query string false "Единоразовый refresh токен, если grant_type=refresh_token"
// @Param       grant_type query string true "Тип разрешения" Enums(authorization_code, refresh_token)
// @Success 200 {object} response.AccessToken "Токен доступа"
// @Failure default {object} response.Error "Ошибка запроса"
// @Router      /oauth/token [post]
func (h *TokenHandler) GenerateToken(c echo.Context) error {
	var err error

	ctx := c.Request().Context()

	clientId := c.QueryParam("client_id")
	clientSecret := c.QueryParam("client_secret")
	codeToken := c.QueryParam("code")
	refreshToken := c.QueryParam("refresh_token")
	grantType := c.QueryParam("grant_type")

	dtoClient := dto.ClientForToken{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	client, err := h.client.ClientForToken(ctx, dtoClient)
	if exception.Is(err) {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	dtoToken := dto.AccessToken{
		Client:    *client,
		Code:      codeToken,
		Refresh:   refreshToken,
		GrantType: grantType,
	}

	access, refresh, err := h.token.AccessAndRefreshToken(ctx, dtoToken)
	if errors.Is(err, exception.ClientAccessDenied) {
		return echo.NewHTTPError(http.StatusForbidden).SetInternal(err)
	}
	if exception.Is(err) {
		return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	userInfo := response.UserInfo{
		UID:   access.User.Id,
		Name:  access.User.Name,
		Email: access.User.Email,
	}

	return c.JSON(http.StatusOK, response.AccessToken{
		TokenType:    "bearer",
		AccessToken:  access.Hash,
		RefreshToken: refresh.Hash,
		ExpiresIn:    access.Expiration.Unix(),
		Info:         userInfo,
	})
}
