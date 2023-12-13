package response

type UserInfo struct {
	// ID пользовтаеля
	UID string `json:"uid"`
	// Полное имя пользовтаеля
	Name string `json:"name"`
	// Email пользователя
	Email string `json:"email"`
}
