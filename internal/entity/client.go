package entity

import "time"

type Client struct {
	Id        string     `db:"id"`
	Name      string     `db:"name"`
	Icon      *string    `db:"icon"`
	Secret    string     `db:"secret"`
	Callback  string     `db:"callback"`
	IsSystem  bool       `db:"is_system"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type ClientRole struct {
	*Client `db:""`
	Role    *string
}
