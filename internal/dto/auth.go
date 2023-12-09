package dto

import "github.com/alnovi/sso/internal/entity"

type AuthById struct {
	Client entity.Client
	UserId string
	IP     string
	Agent  string
}

type AuthByCredentials struct {
	Client   entity.Client
	Login    string
	Password string
	IP       string
	Agent    string
}
