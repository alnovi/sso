package entity

import "time"

type Client struct {
	ID          string
	Name        string
	Description *string
	Icon        *string
	Color       *string
	Image       *string
	Secret      string
	Home        string
	Callback    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
