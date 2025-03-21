package request

type Authorize struct {
	Login    string `json:"login"    form:"login"    validate:"required,email,min=5"        example:"name@example.com"`
	Password string `json:"password" form:"password" validate:"required,gte=5,lte=24" example:"qwerty"`
	Remember bool   `json:"remember" form:"remember"`
}
