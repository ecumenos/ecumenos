package service

import (
	"github.com/ecumenos/ecumenos/pds/config"
	"github.com/ecumenos/ecumenos/pds/repository"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(New),
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
	}
}
