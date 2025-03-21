package response

import (
	"fmt"
	"time"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/pkg/utils"
)

type ProfileUser struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Icon     *string `json:"icon"`
	Callback string  `json:"callback"`
	Role     string  `json:"role"`
}

func NewProfileClient(client *entity.ClientRole) *ProfileClient {
	return &ProfileClient{
		Id:       client.Id,
		Name:     client.Name,
		Icon:     client.Icon,
		Callback: client.Callback,
		Role:     client.Role,
	}
}

func NewCollProfileClient(clients []*entity.ClientRole) []*ProfileClient {
	return utils.MapArray[*ProfileClient, *entity.ClientRole](clients, func(_ int, client *entity.ClientRole) *ProfileClient {
		return NewProfileClient(client)
	})
}

type ProfileSession struct {
	Id        string    `json:"id"`
	IP        string    `json:"ip"`
	App       string    `json:"app"`
	OS        string    `json:"os"`
	IsCurrent bool      `json:"is_current"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
