package entity

const (
	RoleGuestWeight = iota
	RoleUserWeight
	RoleManagerWeight
	RoleAdminWeight

	RoleGuest   = "guest"
	RoleUser    = "user"
	RoleManager = "manager"
	RoleAdmin   = "admin"
)

var RoleMap = map[string]int{
	RoleGuest:   RoleGuestWeight,
	RoleUser:    RoleUserWeight,
	RoleManager: RoleManagerWeight,
	RoleAdmin:   RoleAdminWeight,
}

type Role struct {
	ClientId string `db:"client_id"`
	UserId   string `db:"user_id"`
	Role     string `db:"role"`
}
