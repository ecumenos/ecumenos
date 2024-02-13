package service

import (
	"context"

	"github.com/ecumenos/ecumenos/pds/config"
	"github.com/ecumenos/ecumenos/pds/repository"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
	}
}

type PingServicesResult struct {
	PostgresIsRunning bool
}

func (s *Service) PingServices(ctx context.Context) *PingServicesResult {
	return &PingServicesResult{
		PostgresIsRunning: s.repo.Ping(ctx) == nil,
	}
}
