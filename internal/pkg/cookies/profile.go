package cookies

import "net/http"

const (
	NameProfileAccess  = "pa_hash"
	NameProfileRefresh = "pr_hash"
)

func ProfileAccess(hash string) *http.Cookie {
	return &http.Cookie{
		Name:     NameProfileAccess,
		Value:    hash,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

func ProfileRefresh(hash string) *http.Cookie {
	return &http.Cookie{
		Name:     NameProfileRefresh,
		Value:    hash,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}
