package service

import (
	"context"

	"github.com/ecumenos/ecumenos/orbissocius/config"
	"github.com/ecumenos/ecumenos/orbissocius/repository"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) PingServices(ctx context.Context) *map[string]interface{} {
	return &map[string]interface{}{
		"postgres": s.repo.Ping(ctx) == nil,
	}
}
