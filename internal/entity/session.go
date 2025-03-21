package entity

import (
	"time"

	"github.com/mileusna/useragent"
)

type Session struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	Ip        string    `db:"ip"`
	Agent     string    `db:"agent"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type SessionUser struct {
	*Session
	*User
}

func (s *Session) Parse() useragent.UserAgent {
	return useragent.Parse(s.Agent)
}
