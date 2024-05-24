package entity

import (
	"time"

	"github.com/alnovi/sso/internal/exception"
)

const (
	TokenTypeCode          = "code"
	TokenTypeAccess        = "access"
	TokenTypeRefresh       = "refresh"
	TokenTypeResetPassword = "reset-password"
)

type Token struct {
	ID         string
	Type       string
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
