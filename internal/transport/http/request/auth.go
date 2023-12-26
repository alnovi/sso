package request

// SignInJson request info
// @Description Авторизации пользователя по логину и паролю
type SignInJson struct {
	// Логин или email пользователя
	Login string `json:"login" validate:"required,min=3"`
	// Пароль пользователя
	Password string `json:"password" validate:"required,gte=5,lte=24"`
	// Запомнить меня
	IsRemember bool `json:"remember"`
}

// SignInForm request info
// @Description Авторизации пользователя по логину и паролю
type SignInForm struct {
	// Логин или email пользователя
	Login string `form:"login" validate:"required,min=3"`
	// Пароль пользователя
	Password string `form:"password" validate:"required,gte=5,lte=24"`
	// Запомнить меня
	IsRemember bool `form:"remember"`
}
