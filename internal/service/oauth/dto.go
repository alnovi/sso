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
