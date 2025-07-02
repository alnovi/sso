package response

import (
	"fmt"
	"time"

	"github.com/alnovi/gomon/utils"

	"github.com/alnovi/sso/internal/entity"
)

type ProfileUser struct {
	Id        string    `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name      string    `json:"name" example:"Ivanov"`
	Email     string    `json:"email" example:"ivanov@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2020-01-01T00:00:00+03:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-01-01T00:00:00+03:00"`
}

func NewProfileUser(user *entity.User) *ProfileUser {
	return &ProfileUser{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type ProfileClient struct {
	Id       string  `json:"id" example:"music"`
	Name     string  `json:"name" example:"Music"`
	Icon     *string `json:"icon" example:"icon-music"`
	Callback string  `json:"callback" example:"/callback/auth/"`
	Role     string  `json:"role" example:"admin"`
}

func NewProfileClient(client *entity.ClientRole) *ProfileClient {
	return &ProfileClient{
		Id:       client.Id,
		Name:     client.Name,
		Icon:     client.Icon,
		Callback: client.Callback,
		Role:     *client.Role,
	}
}

func NewCollProfileClient(clients []*entity.ClientRole) []*ProfileClient {
	return utils.MapArray[*ProfileClient, *entity.ClientRole](clients, func(_ int, client *entity.ClientRole) *ProfileClient {
		return NewProfileClient(client)
	})
}

type ProfileSession struct {
	Id        string    `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	IP        string    `json:"ip" example:"127.0.0.1"`
	App       string    `json:"app" example:"music"`
	OS        string    `json:"os" example:"linux"`
	IsCurrent bool      `json:"is_current" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2020-01-01T00:00:00+03:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-01-01T00:00:00+03:00"`
}

func NewProfileSession(session *entity.Session, currentId string) *ProfileSession {
	data := session.Parse()

	return &ProfileSession{
		Id:        session.Id,
		IP:        session.Ip,
		App:       fmt.Sprintf("%s %s", data.Name, data.Version),
		OS:        data.OS,
		IsCurrent: session.Id == currentId,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}
}

func NewCollProfileSession(sessions []*entity.Session, currentId string) []*ProfileSession {
	return utils.MapArray[*ProfileSession, *entity.Session](sessions, func(_ int, session *entity.Session) *ProfileSession {
		return NewProfileSession(session, currentId)
	})
}
