package cookies

import (
	"net/http"
	"time"
)

const (
	NameAuth = "uid"
)

func Auth(userId string, isRemember bool) *http.Cookie {
	cookie := http.Cookie{
		Name:     NameAuth,
		Value:    userId,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	if isRemember {
		cookie.Expires = time.Now().AddDate(0, 1, 0)
	}

	return &cookie
}
