package entity

import "time"

const (
	TokenClassCode    = "code"
	TokenClassAccess  = "access"
	TokenClassRefresh = "refresh"
	TokenClassForgot  = "forgot"

	TokenCodeCost    = 50
	TokenRefreshCost = 100

	TokenCodeTTL    = time.Minute
	TokenAccessTTL  = time.Minute * 2
	TokenRefreshTTL = time.Hour * 24 * 30
)

type Token struct {
	Id         string    `db:"id"`
	Class      string    `db:"class"`
	Hash       string    `db:"hash"`
	SessionId  *string   `db:"session_id"`
	UserId     *string   `db:"user_id"`
	ClientId   *string   `db:"client_id"`
	Payload    Payload   `db:"payload"`
	NotBefore  time.Time `db:"not_before"`
	Expiration time.Time `db:"expiration"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (e *Token) IsActive() bool {
	now := time.Now()

	if now.Before(e.NotBefore) {
		return false
	}

	if now.After(e.Expiration) {
		return false
	}

	return true
}
