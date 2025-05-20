package mailing

import (
	"fmt"
	"strings"
)

type Option func(m *Mailing)

func WithAppHost(host string) Option {
	return func(m *Mailing) {
		m.host = strings.Trim(host, "/")
	}
}

func WithFrom(name, email string) Option {
	return func(m *Mailing) {
		m.from = fmt.Sprintf(`"%s" <%s>`, name, email)
	}
}

func WithAuthUsername(user string) Option {
	return func(m *Mailing) {
		m.client.SetUsername(user)
	}
}

func WithAuthPassword(pass string) Option {
	return func(m *Mailing) {
		m.client.SetPassword(pass)
	}
}
