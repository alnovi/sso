package response

import "time"

type Profile struct {
	Id      string          `json:"uid"`
	Name    string          `json:"name"`
	Image   *string         `json:"image"`
	Email   string          `json:"email"`
	Clients []ProfileClient `json:"clients"`
	Tokens  []ProfileToken  `json:"tokens"`
}

type ProfileClient struct {
	Id          string
	Name        string
	Description *string
	Logo        *string
}

type ProfileToken struct {
	Id        string            `json:"uid"`
	Class     string            `json:"class"`
	Meta      *ProfileTokenMeta `json:"meta"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type ProfileTokenMeta struct {
	IP    string `json:"ip"`
	Agent string `json:"agent"`
}
