package oauth

type InputAuthByCode struct {
	ClientId    string
	RedirectUri string
	Login       string
	Password    string
	State       string
	IP          string
	Agent       string
}

type InputTokenByCode struct {
	ClientId     string
	ClientSecret string
	Code         string
}

type InputTokenByRefresh struct {
	ClientId     string
	ClientSecret string
	Refresh      string
}
