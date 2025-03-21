package cookie

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alnovi/sso/internal/entity"
)

const (
	SessionId    = "session_id"
	accessToken  = "access"
	refreshToken = "refresh"
	sessionTTL   = time.Hour * 24 * 30
)

type Cookie struct {
	secure bool
}

func New(secure bool) *Cookie {
	return &Cookie{secure: secure}
}

func (c *Cookie) SessionId(val string, remember bool) *http.Cookie {
	cookie := &http.Cookie{
		Name:     SessionId,
		Value:    val,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: http.SameSiteStrictMode,
	}

	if remember {
		cookie.Expires = time.Now().Add(sessionTTL)
	}

	return cookie
}

func (c *Cookie) AccessToken(token *entity.Token) *http.Cookie {
	if token == nil || token.Class != entity.TokenClassAccess {
		return nil
	}

	return &http.Cookie{
		Name:     NameAccessToken(*token.ClientId),
		Value:    token.Hash,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: http.SameSiteStrictMode,
		Expires:  token.Expiration,
	}
}

func (c *Cookie) RefreshToken(token *entity.Token) *http.Cookie {
	if token == nil || token.Class != entity.TokenClassRefresh {
		return nil
	}

	return &http.Cookie{
		Name:     NameRefreshToken(*token.ClientId),
		Value:    token.Hash,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: http.SameSiteStrictMode,
		Expires:  token.Expiration,
	}
}

func (c *Cookie) Remove(name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now(),
	}
}

func NameAccessToken(clientId string) string {
	return fmt.Sprintf("%s-%s", clientId, accessToken)
}

func NameRefreshToken(clientId string) string {
	return fmt.Sprintf("%s-%s", clientId, refreshToken)
}
