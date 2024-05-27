package entity

import (
	"time"

	"github.com/alnovi/sso/internal/exception"
)

const (
	TokenClassCode          = "code"
	TokenClassAccess        = "access"
	TokenClassRefresh       = "refresh"
	TokenClassResetPassword = "reset-password"
)

type Token struct {
	ID         string
	Class      string
	Hash       string
	UserID     string
	ClientID   string
	Payload    Payload
	NotBefore  time.Time
	Expiration time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (e *Token) IsActive() error {
	if time.Now().Before(e.NotBefore) {
		return exception.ErrTokenUsedBeforeIssued
	}

	if time.Now().After(e.Expiration) {
		return exception.ErrTokenExpired
	}

	return nil
}
