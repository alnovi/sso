package dto

type ClientForAuth struct {
	ClientId    string
	RedirectURI string
}

type ClientForToken struct {
	ClientId     string
	ClientSecret string
}
