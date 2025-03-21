package request

type CreateUser struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,gte=5,lte=24"`
}

type UpdateUser struct {
	Name     string  `json:"name" validate:"required,min=3,max=100"`
	Email    string  `json:"email" validate:"required,email,max=100"`
	Password *string `json:"password" validate:"omitnil,gte=5,lte=24"`
}

type UpdateUserRole struct {
	Role *string `json:"role" validate:"omitnil,oneof=guest user manager admin"`
}
