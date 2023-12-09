package entity

import "time"

const (
	ClientClassManager = "manager"
	ClientClassProfile = "profile"
	ClientClassClient  = "client"
)

type Client struct {
	Id          string
	Class       string
	Name        string
	Description *string
	Logo        *string
	Image       *string
	Secret      string
	Callback    string
	CanUse      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
