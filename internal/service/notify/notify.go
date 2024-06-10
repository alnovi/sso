package notify

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"

	"github.com/alnovi/sso/internal/entity"
)

const (
	tmplResetPassword = "reset_password.gohtml"
)

//go:embed messages/*.gohtml
var messagesFS embed.FS

type mail interface {
	Sent(ctx context.Context, to, subject, body string) error
}

type Notify struct {
	host string
	mail mail
}

func New(host string, mail mail) *Notify {
	return &Notify{host: host, mail: mail}
}

func (n *Notify) ResetPassword(ctx context.Context, user *entity.User, token *entity.Token) error {
	data := struct {
		UserName   string
		UserEmail  string
		Link       string
		Expiration string
		IP         string
		Agent      string
	}{
		UserName:   user.Name,
		UserEmail:  user.Email,
		Link:       fmt.Sprintf("%s/oauth/reset-password?hash=%s", n.host, token.Hash),
		Expiration: token.Expiration.Format("02.01.2006 15:04"),
		IP:         token.Payload.IP(),
		Agent:      token.Payload.Device(),
	}

	return n.sent(ctx, user.Email, "Востановдение доступа", tmplResetPassword, data)
}

func (n *Notify) sent(ctx context.Context, email, subject, tmplMsg string, data any) error {
	var body bytes.Buffer

	pattern := fmt.Sprintf("messages/%s", tmplMsg)

	tmpl, err := template.ParseFS(messagesFS, pattern, "messages/layout.gohtml")
	if err != nil {
		return err
	}

	err = tmpl.Execute(&body, data)
	if err != nil {
		return err
	}

	return n.mail.Sent(ctx, email, subject, body.String())
}
