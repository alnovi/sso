package server

import (
	"context"

	"github.com/alnovi/sso/internal/entity"
)

type Clients struct {
	manager *entity.Client
	profile *entity.Client
}

func newClients(s *Services) (*Clients, error) {
	manager, err := s.client.GetManagerClient(context.Background())
	if err != nil {
		return nil, err
	}

	profile, err := s.client.GetProfileClient(context.Background())
	if err != nil {
		return nil, err
	}

	return &Clients{
		manager: manager,
		profile: profile,
	}, nil
}
