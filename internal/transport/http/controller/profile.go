package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/cookie"
	"github.com/alnovi/sso/internal/service/profile"
	"github.com/alnovi/sso/internal/transport/http/request"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/pkg/validator"
)

type ProfileController struct {
	BaseController
	profile *profile.UserProfile
	cookie  *cookie.Cookie
	session echo.MiddlewareFunc
}

func NewProfileController(profile *profile.UserProfile, cookie *cookie.Cookie, session echo.MiddlewareFunc) *ProfileController {
	return &ProfileController{profile: profile, cookie: cookie, session: session}
}

func (c *ProfileController) Home(e echo.Context) error {
	sessionId := e.QueryParam(cookie.SessionId)
	userAgent := e.Request().UserAgent()

	if _, err := e.Cookie(cookie.SessionId); err != nil {
		if _, err = c.profile.SessionByIdAndAgent(context.Background(), sessionId, userAgent); err == nil {
			e.SetCookie(c.cookie.SessionId(sessionId, false))
		}
	}

	if sessionId != "" {
		return e.Redirect(http.StatusFound, "/profile")
	}

	return e.Render(http.StatusOK, "profile.html", nil)
}

func (c *ProfileController) Me(e echo.Context) error {
	userId := c.MustUserId(e)

	user, err := c.profile.Info(context.Background(), userId)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, response.NewProfileUser(user))
}

func (c *ProfileController) UpdateUser(e echo.Context) error {
	userId := c.MustUserId(e)

	req := new(request.UpdateProfile)

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	user, err := c.profile.UpdateInfo(context.Background(), userId, req.Name, req.Email)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, response.NewProfileUser(user))
}

func (c *ProfileController) Clients(e echo.Context) error {
	userId := c.MustUserId(e)

	clients, err := c.profile.Clients(context.Background(), userId)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, response.NewCollProfileClient(clients))
}

func (c *ProfileController) Sessions(e echo.Context) error {
	userId := c.MustUserId(e)
	sessionId := c.MustSessionId(e)

	sessions, err := c.profile.Sessions(context.Background(), userId)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, response.NewCollProfileSession(sessions, sessionId))
}

func (c *ProfileController) SessionDelete(e echo.Context) error {
	userId := c.MustUserId(e)
	sessionId := c.MustSessionId(e)
	id := e.Param("id")

	err := c.profile.SessionDelete(context.Background(), userId, id)
	if err != nil {
		if errors.Is(err, profile.ErrSessionNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "session not found").SetInternal(err)
		}
		return err
	}

	if sessionId == id {
		e.SetCookie(c.cookie.Remove(cookie.SessionId))
	}

	return e.NoContent(http.StatusOK)
}

func (c *ProfileController) UpdatePassword(e echo.Context) error {
	userId := c.MustUserId(e)

	req := new(request.UpdatePassword)

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	err := c.profile.UpdatePassword(context.Background(), userId, req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, profile.ErrInvalidPassword) {
			return validator.NewValidateErrorWithMessage("old_password", "Пароль не верный")
		}
		return err
	}

	return e.NoContent(http.StatusOK)
}

func (c *ProfileController) Logout(e echo.Context) error {
	_ = c.profile.Logout(context.Background(), c.MustSessionId(e))
	e.SetCookie(c.cookie.Remove(cookie.SessionId))
	return e.NoContent(http.StatusOK)
}

func (c *ProfileController) ApplyHTTP(g *echo.Group) error {
	g.GET("/profile/", c.Home, c.session)
	g.GET("/profile/me/", c.Me, c.session)
	g.PUT("/profile/me/", c.UpdateUser, c.session)
	g.GET("/profile/clients/", c.Clients, c.session)
	g.GET("/profile/sessions/", c.Sessions, c.session)
	g.DELETE("/profile/sessions/:id/", c.SessionDelete, c.session)
	g.PUT("/profile/password/", c.UpdatePassword, c.session)
	g.POST("/profile/logout/", c.Logout, c.session)
	return nil
}
