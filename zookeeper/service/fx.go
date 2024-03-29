package service

import (
	"github.com/ecumenos/ecumenos/internal/fxappsettings"
	"github.com/ecumenos/ecumenos/zookeeper/config"
	"github.com/ecumenos/ecumenos/zookeeper/repository"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(New),
)

type Service struct {
	repo     *repository.Repository
	auth     *Authorization
	settings fxappsettings.AppSettings
}

func New(repo *repository.Repository, rm fxappsettings.AppSettings, cfg *config.Config) *Service {
	return &Service{
		repo:     repo,
		auth:     &Authorization{JWTSigningKey: cfg.JWTSecret},
		settings: rm,
	}
}
