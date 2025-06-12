package response

import (
	"time"

	"github.com/alnovi/gomon/utils"

	"github.com/alnovi/sso/internal/entity"
)

type User struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func NewUser(user *entity.User) *User {
	return &User{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

func NewUsers(users []*entity.User) []*User {
	return utils.MapArray[*User, *entity.User](users, func(_ int, user *entity.User) *User {
		return NewUser(user)
	})
}
