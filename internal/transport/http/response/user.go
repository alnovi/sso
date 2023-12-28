package response

type User struct {
	UID   string  `json:"uid"`
	Name  string  `json:"name"`
	Image *string `json:"image"`
	Email string  `json:"email"`
}
