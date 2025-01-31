package response

import "time"

type Client struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AccessToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in"`
}
