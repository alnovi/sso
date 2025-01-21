package cookie

import (
	"net/http"
	"time"
)

const (
	SessionId = "session_id"
)

type Cookie struct {
	secure bool
}

func New(secure bool) *Cookie {
	return &Cookie{secure: secure}
}

func (c *Cookie) SessionId(val string) *http.Cookie {
	return &http.Cookie{
		Name:     SessionId,
		Value:    val,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.secure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Hour * 24 * 30), //nolint:mnd
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
