package response

import (
	"time"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/pkg/utils"
)

type Client struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Icon      *string    `json:"icon"`
	Secret    string     `json:"secret"`
	Callback  string     `json:"callback"`
	IsSystem  bool       `json:"is_system"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func NewClient(client *entity.Client) *Client {
	return &Client{
		Id:        client.Id,
		Name:      client.Name,
		Icon:      client.Icon,
		Secret:    client.Secret,
		Callback:  client.Callback,
		IsSystem:  client.IsSystem,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
		DeletedAt: client.DeletedAt,
	}
}

func NewClients(clients []*entity.Client) []*Client {
	return utils.MapArray[*Client, *entity.Client](clients, func(_ int, client *entity.Client) *Client {
		return NewClient(client)
	})
}

type ClientRole struct {
	*Client
	Role string `json:"role"`
}

func NewClientRole(client *entity.ClientRole) *ClientRole {
	return &ClientRole{
		Client: NewClient(client.Client),
		Role:   client.Role,
	}
}

func NewClientsRoles(clients []*entity.ClientRole) []*ClientRole {
	return utils.MapArray[*ClientRole, *entity.ClientRole](clients, func(_ int, client *entity.ClientRole) *ClientRole {
		return NewClientRole(client)
	})
}
