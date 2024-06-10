package server

import (
	"github.com/alnovi/sso/internal/service/notify"
	"github.com/alnovi/sso/internal/service/secure"
)

func (p *Provider) Notify() *notify.Notify {
	if p.notify == nil {
		p.notify = notify.New(p.Config().App.Host, p.Mail())
	}
	return p.notify
}

func (p *Provider) Secure() *secure.Secure {
	if p.secure == nil {
		p.secure = secure.NewSecure(p.Repository())
	}
	return p.secure
}
