package response

// AccessToken response info
// @Description Токен доступа
type AccessToken struct {
	// Тип авторизации
	TokenType string `json:"token_type"`
	// Токен доступа
	AccessToken string `json:"access_token"`
	// Токен обновления
	RefreshToken string `json:"refresh_token"`
	// Время действия токена
	ExpiresIn int64 `json:"expires_in"`
	// Информация о пользователе
	Info UserInfo `json:"info"`
}
