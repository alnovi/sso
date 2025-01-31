package request

type Authorize struct {
	Login    string `json:"login"    form:"login"    validate:"required,min=5"        example:"name@example.com"`
	Password string `json:"password" form:"password" validate:"required,gte=5,lte=24" example:"qwerty"`
	Remember bool   `json:"remember" form:"remember"`
}

type TokenByCode struct {
	ClientId     string `query:"client_id" validate:"required"`
	ClientSecret string `query:"client_secret" validate:"required"`
	Code         string `query:"code" validate:"required"`
}

type TokenByRefresh struct {
	ClientId     string `query:"client_id" validate:"required"`
	ClientSecret string `query:"client_secret" validate:"required"`
	Refresh      string `query:"refresh_token" validate:"required"`
}

type ForgotPassword struct {
	Login string `json:"login" form:"login" validate:"required,min=5" example:"name@example.com"`
}

type ResetPassword struct {
	Password string `json:"password" form:"password" validate:"required,gte=5,lte=24" example:"qwerty"`
}
