package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service/sessions"
	"github.com/alnovi/sso/internal/service/token"
	"github.com/alnovi/sso/pkg/utils"
)

const (
	ResponseTypeCode           = "code"
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypeRefreshToken      = "refresh_token"
)

var (
	ErrClientNotFound      = errors.New("client not found")
	ErrTokenNotFound       = errors.New("token not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserPasswordInvalid = errors.New("user password invalid")
	ErrResponseTypeInvalid = errors.New("response type invalid")
	ErrRedirectUriInvalid  = errors.New("redirect uri invalid")
)

type OAuth struct {
	repo    repository.Repository
	tm      repository.Transaction
	token   *token.Token
	session *sessions.Session
}

func New(r repository.Repository, tm repository.Transaction, t *token.Token, s *sessions.Session) *OAuth {
	return &OAuth{repo: r, tm: tm, token: t, session: s}
}

func (s *OAuth) Client(ctx context.Context, id string, secret *string) (*entity.Client, error) {
	client, err := s.repo.ClientById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w [client_id=%s]: %s", ErrClientNotFound, id, err)
	}

	if !client.IsActive {
		return nil, fmt.Errorf("%w [client_id=%s]: client is not active", ErrClientNotFound, id)
	}

	if secret != nil && client.Secret != *secret {
		return nil, fmt.Errorf("%w [client_id=%s]: secret does not match", ErrClientNotFound, id)
	}

	return client, nil
}

func (s *OAuth) ResponseType(rt string) (string, error) {
	rt = strings.ToLower(rt)

	allowedTypes := []string{ResponseTypeCode}

	if !slices.Contains(allowedTypes, rt) {
		return rt, ErrResponseTypeInvalid
	}

	return rt, nil
}

func (s *OAuth) RedirectURL(client *entity.Client, uri string) (*url.URL, error) {
	if client == nil {
		return nil, ErrClientNotFound
	}

	clientURI, err := url.Parse(client.Host)
	if err != nil {
		return nil, fmt.Errorf("%w: can't parse client host [client_id=%s]: %s", ErrRedirectUriInvalid, client.Id, err)
	}

	redirectURI, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("%w: can't parse redirect_uri [uri=%s]: %s", ErrRedirectUriInvalid, uri, err)
	}

	if clientURI.Host != redirectURI.Host {
		return nil, fmt.Errorf("%w: redirect_uri host does not match [client_id=%s] [uri=%s]", ErrRedirectUriInvalid, client.Id, uri)
	}

	return redirectURI, nil
}

func (s *OAuth) AuthorizeByCode(ctx context.Context, inp InputAuthByCode) (*url.URL, *entity.Token, error) {
	var code *entity.Token

	client, err := s.Client(ctx, inp.ClientId, nil)
	if err != nil {
		return nil, nil, err
	}

	redirectUri, err := s.RedirectURL(client, inp.RedirectUri)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.userAttempt(ctx, inp.Login, inp.Password)
	if err != nil {
		return nil, nil, fmt.Errorf("%w [user_email=%s]", err, inp.Login)
	}

	err = s.tm.ReadCommitted(ctx, func(ctx context.Context) error {
		var session *entity.Session

		session, err = s.session.Create(ctx, user.Id, inp.IP, inp.Agent)
		if err != nil {
			return fmt.Errorf("fail create session: %s", err)
		}

		code, err = s.token.CodeToken(ctx, session.Id, client.Id, user.Id)
		if err != nil {
			return fmt.Errorf("can't create code token [session_id=%s] [client_id=%s] [user_id=%s]: %w", session.Id, client.Id, user.Id, err)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	query := redirectUri.Query()
	query.Set("code", code.Hash)
	query.Set("state", inp.State)
	redirectUri.RawQuery = query.Encode()

	return redirectUri, code, nil
}

func (s *OAuth) TokenByCode(ctx context.Context, inp InputTokenByCode) (*entity.Token, *entity.Token, error) {
	var accessToken *entity.Token
	var refreshToken *entity.Token

	client, err := s.Client(ctx, inp.ClientId, &inp.ClientSecret)
	if err != nil {
		return nil, nil, err
	}

	err = s.tm.ReadCommitted(ctx, func(ctx context.Context) error {
		var code *entity.Token
		var role *entity.Role

		code, err = s.repo.TokenByClassHash(ctx, entity.TokenClassCode, inp.Code)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrTokenNotFound, err)
		}

		if !code.CheckClient(client.Id) {
			return fmt.Errorf("%w: client is not match", ErrTokenNotFound)
		}

		if !code.IsActive() {
			return fmt.Errorf("%w: code is not active", ErrTokenNotFound)
		}

		if err = s.repo.TokenDelete(ctx, code.Id); err != nil {
			return fmt.Errorf("can't delete code-token: %s", err)
		}

		role, err = s.repo.RoleByClientAndUser(ctx, *code.ClientId, *code.UserId)
		if err != nil {
			return fmt.Errorf("role not found: %s", err)
		}

		accessToken, err = s.token.AccessToken(ctx, *code.SessionId, client.Id, *code.UserId, role.Role)
		if err != nil {
			return err
		}

		refreshToken, err = s.token.RefreshToken(ctx, *code.SessionId, client.Id, *code.UserId, accessToken.Expiration)
		if err != nil {
			return err
		}

		return nil
	})

	return accessToken, refreshToken, err
}

func (s *OAuth) TokenByRefresh(ctx context.Context, inp InputTokenByRefresh) (*entity.Token, *entity.Token, error) {
	var accessToken *entity.Token
	var refreshToken *entity.Token

	client, err := s.Client(ctx, inp.ClientId, &inp.ClientSecret)
	if err != nil {
		return nil, nil, err
	}

	err = s.tm.ReadCommitted(ctx, func(ctx context.Context) error {
		var refresh *entity.Token
		var role *entity.Role

		refresh, err = s.repo.TokenByClassHash(ctx, entity.TokenClassRefresh, inp.Refresh)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrTokenNotFound, err)
		}

		if !refresh.CheckClient(client.Id) {
			return fmt.Errorf("%w: client is not match", ErrTokenNotFound)
		}

		if !refresh.IsActive() {
			return fmt.Errorf("%w: refresh is not active", ErrTokenNotFound)
		}

		if err = s.repo.TokenDelete(ctx, refresh.Id); err != nil {
			return fmt.Errorf("can't delete refresh-token: %s", err)
		}

		role, err = s.repo.RoleByClientAndUser(ctx, *refresh.ClientId, *refresh.UserId)
		if err != nil {
			return fmt.Errorf("role not found: %s", err)
		}

		accessToken, err = s.token.AccessToken(ctx, *refresh.SessionId, client.Id, *refresh.UserId, role.Role)
		if err != nil {
			return err
		}

		refreshToken, err = s.token.RefreshToken(ctx, *refresh.SessionId, client.Id, *refresh.UserId, accessToken.Expiration)
		if err != nil {
			return err
		}

		return nil
	})

	return accessToken, refreshToken, err
}

func (s *OAuth) RemoveSession(ctx context.Context, sessionId string) error {
	return s.session.Delete(ctx, sessionId)
}

func (s *OAuth) userAttempt(ctx context.Context, email, password string) (*entity.User, error) {
	user, err := s.repo.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNoResults) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if !utils.CompareHashPassword(password, user.Password) {
		return nil, ErrUserPasswordInvalid
	}

	return user, nil
}
