package mail

import (
	"context"
	"fmt"

	"github.com/wneessen/go-mail"
)

type From struct {
	Name  string
	Email string
}

type Mail struct {
	client *mail.Client
	from   string
}

func New(from From, user, password, host string, port int) (*Mail, error) {
	client, err := mail.NewClient(
		host,
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(user),
		mail.WithPassword(password),
		mail.WithTLSPortPolicy(mail.TLSOpportunistic),
		mail.WithPort(port),
	)

	return &Mail{
		client: client,
		from:   fmt.Sprintf(`"%s" <%s>`, from.Name, from.Email),
	}, err
}

func (m *Mail) Sent(ctx context.Context, to, subject, body string) error {
	var err error

	msg := mail.NewMsg()

	if err = msg.From(m.from); err != nil {
		return err
	}

	if err = msg.To(to); err != nil {
		return err
	}

	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextHTML, body)

	return m.client.DialAndSendWithContext(ctx, msg)
}
