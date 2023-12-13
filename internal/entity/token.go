package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

const (
	TokenClassCode    = "code"
	TokenClassAccess  = "access"
	TokenClassRefresh = "refresh"
)

type Token struct {
	Id         string
	Class      string
	Hash       string
	UserId     *string
	ClientId   *string
	Meta       *TokenMeta
	NotBefore  time.Time
	Expiration time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TokenMeta struct {
	IP    string `json:"ip,omitempty"`
	Agent string `json:"agent,omitempty"`
}

type TokenWithUser struct {
	Token
	User User
}

func (tm *TokenMeta) Value() (driver.Value, error) {
	if tm == nil {
		return nil, nil
	}

	return json.Marshal(tm)
}

func (tm *TokenMeta) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, tm)
}
