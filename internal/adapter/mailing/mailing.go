package mailing

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"strconv"

	"github.com/wneessen/go-mail"
	"go.opentelemetry.io/otel/attribute"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/helper"
)

//go:embed messages/*.html
var messagesFS embed.FS

type Mailing struct {
	host   string
	from   string
	client *mail.Client
}

func New(host, port string, opts ...Option) (*Mailing, error) {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("fail convert mailing port: %s", err)
	}

	client, err := mail.NewClient(host,
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTLSPortPolicy(mail.TLSOpportunistic),
		mail.WithPort(portNum),
	)

	if err != nil {
		return nil, fmt.Errorf("create mailing client: %w", err)
	}

	mailing := &Mailing{from: "SSO", client: client}

	for _, opt := range opts {
		opt(mailing)
	}

	return mailing, nil
}

func (m *Mailing) Ping(ctx context.Context) error {
	if err := m.client.DialWithContext(ctx); err != nil {
		return fmt.Errorf("mailing dial fail: %s", err)
	}
	return nil
}

func (m *Mailing) Close(_ context.Context) error {
	return m.client.Close()
}

func (m *Mailing) ForgotPassword(ctx context.Context, user *entity.User, token *entity.Token) error {
	ctx, span := helper.SpanStart(ctx, "Mailing.ForgotPassword", helper.SpanAttr(
		attribute.String("user.id", user.Id),
		attribute.String("user.email", user.Email),
	))
	defer span.End()

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
		Link:       fmt.Sprintf("%s/oauth/reset-password?hash=%s", m.host, token.Hash),
		Expiration: token.Expiration.Format("02.01.2006 15:04"),
		IP:         token.Payload.IP(),
		Agent:      token.Payload.Agent(),
	}

	err := m.sentMsg(ctx, user.Email, "Восстановление доступа", "forgot_password.html", data)
	if err != nil {
		helper.SpanError(span, err)
	}

	return err
}

func (m *Mailing) sentMsg(ctx context.Context, email, subject, tmpl string, data any) error {
	var body bytes.Buffer

	tmpl = fmt.Sprintf("messages/%s", tmpl)

	tmplMsg, err := template.ParseFS(messagesFS, tmpl, "messages/layout.html")
	if err != nil {
		return err
	}

	err = tmplMsg.Execute(&body, data)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()

	if err = msg.From(m.from); err != nil {
		return err
	}

	if err = msg.To(email); err != nil {
		return err
	}

	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextHTML, body.String())

	return m.client.DialAndSendWithContext(ctx, msg)
}
