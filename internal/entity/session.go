package entity

import "time"

type Session struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	Ip        string    `db:"ip"`
	Agent     string    `db:"agent"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
