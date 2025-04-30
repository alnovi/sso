package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"

	"github.com/google/uuid"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service/token"
	"github.com/alnovi/sso/pkg/utils"
)

const (
	ResponseTypeCode           = "code"
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypeRefreshToken      = "refresh_token"
)

var (
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrClientNotFound      = errors.New("client not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrTokenNotFound       = errors.New("token not found")
	ErrSessionNotFound     = errors.New("session not found")
	ErrInvalidUserPassword = errors.New("invalid user password")
	ErrInvalidResponseType = errors.New("invalid response type")
	ErrInvalidRedirectUri  = errors.New("invalid redirect uri")

	responseTypes = []string{ResponseTypeCode}
)

type OAuth struct {
	repo  *repository.Repository
	tm    repository.Transaction
	token *token.Token
}

func NewOAuth(repo *repository.Repository, tm repository.Transaction, token *token.Token) *OAuth {
	return &OAuth{repo: repo, tm: tm, token: token}
}

func (s *OAuth) AuthorizeCheckParams(ctx context.Context, inp InputAuthorizeParams) (*entity.Client, error) {
	if !slices.Contains(responseTypes, inp.ResponseType) {
		return nil, ErrInvalidResponseType
	}

	client, err := s.repo.ClientById(ctx, inp.ClientId)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrClientNotFound, err)
	}

	err = utils.CompareHosts(inp.RedirectUri, client.Callback)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRedirectUri, err)
	}

	return client, nil
}

func (s *OAuth) AuthorizeByCode(ctx context.Context, inp InputAuthorizeByCode) (*entity.Client, *entity.Token, *url.URL, error) {
	var session *entity.Session
	var code *entity.Token

	if !slices.Contains(responseTypes, inp.ResponseType) {
		return nil, nil, nil, ErrInvalidResponseType
	}

	client, err := s.repo.ClientById(ctx, inp.ClientId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %s", ErrClientNotFound, err)
	}

	err = utils.CompareHosts(inp.RedirectUri, client.Callback)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %s", ErrInvalidRedirectUri, err)
	}

	user, err := s.repo.UserByEmail(ctx, inp.Login, repository.NotDeleted())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %s", ErrUserNotFound, err)
	}

	if !utils.CompareHashPassword(inp.Password, user.Password) {
		return nil, nil, nil, ErrInvalidUserPassword
	}

	_, err = s.repo.Role(ctx, client.Id, user.Id)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %s", ErrForbidden, err)
	}

	err = s.tm.ReadCommitted(ctx, func(ctx context.Context) error {
		session, err = s.repo.SessionByUserId(ctx, user.Id, repository.IP(inp.UserIP), repository.Agent(inp.UserAgent))
		if err != nil {
			session = &entity.Session{
				Id:     uuid.NewString(),
				UserId: user.Id,
				Ip:     inp.UserIP,
				Agent:  inp.UserAgent,
			}

			if err = s.repo.SessionCreate(ctx, session); err != nil {
				return err
			}
		}

		code, err = s.token.CodeToken(ctx, session.Id, client.Id, user.Id)

		return err
	})

	if err != nil {
		return nil, nil, nil, err
	}

	redirectUri, _ := url.Parse(inp.RedirectUri)
	query := redirectUri.Query()
	query.Add("code", code.Hash)
	query.Add("state", inp.State)
	redirectUri.RawQuery = query.Encode()

	return client, code, redirectUri, nil
}

func (s *OAuth) AuthorizeBySession(ctx context.Context, inp InputAuthorizeBySession) (*entity.Client, *entity.Token, *url.URL, error) {
	if !slices.Contains(responseTypes, inp.ResponseType) {
		return nil, nil, nil, ErrInvalidResponseType
	}

	client, err := s.repo.ClientById(ctx, inp.ClientId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %s", ErrClientNotFound, err)
	}

	err = utils.CompareHosts(inp.RedirectUri, client.Callback)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %s", ErrInvalidRedirectUri, err)
	}

	session, err := s.repo.SessionById(ctx, inp.SessionId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %s", ErrSessionNotFound, err)
	}

	code, err := s.token.CodeToken(ctx, session.Id, client.Id, session.UserId)
	if err != nil {
		return nil, nil, nil, err
	}

	redirectUri, _ := url.Parse(inp.RedirectUri)
	query := redirectUri.Query()
	query.Add("code", code.Hash)
	query.Add("state", inp.State)
	redirectUri.RawQuery = query.Encode()

	return client, code, redirectUri, nil
}

func (s *OAuth) TokenByCode(ctx context.Context, inp InputTokenByCode) (*entity.Token, *entity.Token, error) {
	var code *entity.Token
	var accessToken *entity.Token
	var refreshToken *entity.Token

	client, err := s.repo.ClientById(ctx, inp.ClientId, repository.Secret(inp.ClientSecret), repository.NotDeleted())
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrClientNotFound, err)
	}

	err = s.tm.ReadCommitted(ctx, func(ctx context.Context) error {
		var role *entity.Role

		code, err = s.repo.TokenByHash(ctx, inp.Code, repository.Class(entity.TokenClassCode), repository.ForUpdate())
		if err != nil {
			return fmt.Errorf("%w: %s", ErrTokenNotFound, err)
		}

		if !code.IsActive() {
			return ErrTokenNotFound
		}

		if code.SessionId == nil {
			return ErrSessionNotFound
		}

		if *code.ClientId != client.Id {
			return ErrTokenNotFound
		}

		if err = s.repo.SessionUpdateDateById(ctx, *code.SessionId); err != nil {
			return err
		}

		if err = s.repo.TokenDeleteBySessionId(ctx, *code.SessionId); err != nil {
			return err
		}

		role, err = s.repo.Role(ctx, client.Id, *code.UserId)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrForbidden, err)
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
	var refresh *entity.Token
	var accessToken *entity.Token
	var refreshToken *entity.Token

	client, err := s.repo.ClientById(ctx, inp.ClientId, repository.Secret(inp.ClientSecret), repository.NotDeleted())
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrClientNotFound, err)
	}

	err = s.tm.ReadCommitted(ctx, func(ctx context.Context) error {
		var role *entity.Role

		refresh, err = s.repo.TokenByHash(ctx, inp.Refresh, repository.Class(entity.TokenClassRefresh), repository.ForUpdate())
		if err != nil {
			return fmt.Errorf("%w: %s", ErrTokenNotFound, err)
		}

		if !refresh.IsActive() {
			return ErrTokenNotFound
		}

		if refresh.SessionId == nil {
			return ErrSessionNotFound
		}

		if *refresh.ClientId != client.Id {
			return ErrTokenNotFound
		}

		if err = s.repo.SessionUpdateDateById(ctx, *refresh.SessionId); err != nil {
			return err
		}

		if err = s.repo.TokenDeleteById(ctx, refresh.Id); err != nil {
			return err
		}

		role, err = s.repo.Role(ctx, client.Id, *refresh.UserId)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrForbidden, err)
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

func (s *OAuth) ValidateAccessToken(ctx context.Context, token string) (*token.AccessClaims, error) {
	claims, err := s.token.ValidateAccessToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnauthorized, err)
	}
	return claims, nil
}

func (s *OAuth) ValidateRefreshToken(ctx context.Context, token string) (*entity.Token, error) {
	return s.token.ValidateRefreshToken(ctx, token)
}
