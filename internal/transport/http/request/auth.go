package request

type SignIn struct {
	Login      string `json:"login" form:"login" validate:"required,min=3"`
	Password   string `json:"password" form:"password" validate:"required,gte=6,lte=24"`
	IsRemember bool   `json:"remember" form:"remember"`
}
