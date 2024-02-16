package service

import (
	"github.com/ecumenos/ecumenos/orbissocius/config"
	"github.com/ecumenos/ecumenos/orbissocius/repository"
	"go.uber.org/fx"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
	}
}

var Module = fx.Options(
	fx.Provide(New),
)
