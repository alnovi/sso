package server

import "github.com/alnovi/sso/internal/service/secure"

func (p *Provider) Secure() *secure.Secure {
	if p.secure == nil {
		p.secure = secure.NewSecure(p.Repository())
	}
	return p.secure
}
