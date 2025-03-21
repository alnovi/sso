package response

import (
	"fmt"
	"time"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/pkg/utils"
)

type Session struct {
	Id        string    `json:"id"`
	IP        string    `json:"ip"`
	App       string    `json:"app"`
	OS        string    `json:"os"`
	Agent     string    `json:"agent"`
	IsCurrent bool      `json:"is_current"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewSession(session *entity.Session, sessionId string) *Session {
	data := session.Parse()

	return &Session{
		Id:        session.Id,
		IP:        session.Ip,
		App:       fmt.Sprintf("%s %s", data.Name, data.Version),
		OS:        data.OS,
		Agent:     session.Agent,
		IsCurrent: session.Id == sessionId,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}
}

type SessionUser struct {
	*Session
	User *User `json:"user"`
}

func NewSessionUser(session *entity.SessionUser, sessionId string) *SessionUser {
	return &SessionUser{
		Session: NewSession(session.Session, sessionId),
		User:    NewUser(session.User),
	}
}

func NewSessionsUser(sessions []*entity.SessionUser, sessionId string) []*SessionUser {
	return utils.MapArray[*SessionUser, *entity.SessionUser](sessions, func(_ int, item *entity.SessionUser) *SessionUser {
		return NewSessionUser(item, sessionId)
	})
}
