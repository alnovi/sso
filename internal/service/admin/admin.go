package admin

import (
	"context"
	"net/url"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
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
	client, err := s.repo.ClientById(ctx, s.clientId)
	if err != nil {
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
	client, err := s.repo.ClientById(ctx, s.clientId)
	if err != nil {
		return nil, nil, err
	}

	inp := oauth.InputTokenByCode{
		ClientId:     client.Id,
		ClientSecret: client.Secret,
		Code:         code,
	}

	return s.oauth.TokenByCode(ctx, inp)
}

func (s *Admin) Logout(ctx context.Context, sessionId string) error {
	return s.repo.SessionDeleteById(ctx, sessionId)
}
