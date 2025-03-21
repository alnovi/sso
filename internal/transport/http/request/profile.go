package request

type UpdateProfile struct {
	Name  string `json:"name" validate:"required,gte=3,lte=100" example:"Ivanov"`
	Email string `json:"email" validate:"required,email,gte=5,lte=100" example:"ivanov@example.com"`
}

type UpdatePassword struct {
	OldPassword string `json:"old_password" validate:"required,gte=5,lte=24" example:"secret"`
	NewPassword string `json:"new_password" validate:"required,gte=5,lte=24" example:"secret"`
}
