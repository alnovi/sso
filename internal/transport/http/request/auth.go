package request

// SignIn model info
// @Description Авторизации пользователя по логину и паролю
type SignIn struct {
	// Логин или email пользователя
	Login string `json:"login" form:"login" validate:"required,min=3"`
	// Пароль пользователя
	Password string `json:"password" form:"password" validate:"required,gte=6,lte=24"`
	// Запомнить меня
	IsRemember bool `json:"remember" form:"remember"`
}
