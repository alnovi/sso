package entity

import "time"

const (
	TokenClassCode    = "code"
	TokenClassAccess  = "access"
	TokenClassRefresh = "refresh"
	TokenClassForgot  = "forgot"
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

func NewToken() *Token {
	return &Token{Payload: Payload{}}
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

func (e *Token) CheckClient(clientId string) bool {
	return e.ClientId != nil && *e.ClientId == clientId
}
