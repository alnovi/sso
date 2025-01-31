package entity

const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleUser    = "user"
	RoleGuest   = "guest"
)

type Role struct {
	ClientId string `db:"client_id"`
	UserId   string `db:"user_id"`
	Role     string `db:"role"`
}
