package entity

import "time"

type Client struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	Secret    string    `db:"secret"`
	Host      string    `db:"host"`
	Icon      *string   `db:"icon"`
	Color     *string   `db:"color"`
	Image     *string   `db:"image"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
