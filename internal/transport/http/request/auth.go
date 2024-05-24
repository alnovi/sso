package request

type Authorize struct {
	Login      string `form:"login"    json:"login"    validate:"required,min=3,email"`
	Password   string `form:"password" json:"password" validate:"required,gte=5,lte=24"`
	IsRemember bool   `form:"remember" json:"remember"`
}
