package entity

import "time"

type User struct {
	Id        string
	Image     *string
	Name      string
	Login     string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
