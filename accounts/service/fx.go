package service

import (
	"github.com/ecumenos/ecumenos/accounts/config"
	"go.uber.org/fx"
)

type Service struct {
}

func New(cfg *config.Config) *Service {
	return &Service{}
}

var Module = fx.Options(
	fx.Provide(New),
)
