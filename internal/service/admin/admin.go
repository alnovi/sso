package admin

import (
	"context"
	"net/url"

	"go.opentelemetry.io/otel/attribute"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/helper"
	"github.com/alnovi/sso/internal/service/oauth"
)

type Admin struct {
	clientId string
	repo     *repository.Repository
	tm       repository.Transaction
	oauth    *oauth.OAuth
}

func NewAdmin(clientId string, repo *repository.Repository, tm repository.Transaction, oauth *oauth.OAuth) *Admin {
	return &Admin{clientId: clientId, repo: repo, tm: tm, oauth: oauth}
}

func (s *Admin) ClientId() string {
	return s.clientId
}

func (s *Admin) AuthorizeURI(ctx context.Context) (string, error) {
	ctx, span := helper.SpanStart(ctx, "Admin.AuthorizeURI")
	defer span.End()

	client, err := s.repo.ClientById(ctx, s.clientId)
	if err != nil {
		helper.SpanError(span, err)
		return "", err
	}

	query := url.Values{
		"response_type": []string{oauth.ResponseTypeCode},
		"client_id":     []string{client.Id},
		"redirect_uri":  []string{client.Callback},
	}

	authorizeURI := &url.URL{
		Path:     "/oauth/authorize",
		RawQuery: query.Encode(),
	}

	return authorizeURI.String(), nil
}

func (s *Admin) TokenByCode(ctx context.Context, code string) (*entity.Token, *entity.Token, error) {
	ctx, span := helper.SpanStart(ctx, "Admin.TokenByCode", helper.SpanAttr(
		attribute.String("code", code),
	))
	defer span.End()

	client, err := s.repo.ClientById(ctx, s.clientId)
	if err != nil {
		helper.SpanError(span, err)
		return nil, nil, err
	}

	inp := oauth.InputTokenByCode{
		ClientId:     client.Id,
		ClientSecret: client.Secret,
		Code:         code,
	}

	access, refresh, err := s.oauth.TokenByCode(ctx, inp)
	if err != nil {
		helper.SpanError(span, err)
		return nil, nil, err
	}

	return access, refresh, nil
}

func (s *Admin) Logout(ctx context.Context, sessionId string) error {
	ctx, span := helper.SpanStart(ctx, "Admin.Logout", helper.SpanAttr(
		attribute.String("session.id", sessionId),
	))
	defer span.End()

	err := s.repo.SessionDeleteById(ctx, sessionId)
	helper.SpanError(span, err)

	return err
}
