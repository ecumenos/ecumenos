package repository

import (
	"context"

	"github.com/ecumenos/ecumenos/internal/fxpostgres/postgres"
	"github.com/ecumenos/ecumenos/pds/config"
	"go.uber.org/zap"
)

type Repository struct {
	driver *postgres.Driver
	logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) (*Repository, error) {
	driver, err := postgres.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	return &Repository{
		driver: driver,
		logger: logger,
	}, nil
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.driver.Ping(ctx)
}
