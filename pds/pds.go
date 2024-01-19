package pds

import (
	"context"

	"github.com/ecumenos/ecumenos/internal/fxpostgres/postgres"
	"github.com/ecumenos/ecumenos/internal/fxtypes"
	"go.uber.org/zap"
)

var (
	ServiceName    fxtypes.ServiceName = "pds"
	ServiceVersion fxtypes.Version     = "v0.0.0"
)

type Config struct {
	Addr        string
	Prod        bool
	PostgresURL string
}

type PDS struct {
	Postgres *postgres.Driver
	logger   *zap.Logger
}

func New(cfg *Config, l *zap.Logger) (*PDS, error) {
	driver, err := postgres.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	return &PDS{
		Postgres: driver,
		logger:   l,
	}, nil
}

func (s *PDS) Start(ctx context.Context) error {
	if err := s.Postgres.Ping(ctx); err != nil {
		return err
	}
	s.logger.Info("postgres is started")

	return nil
}

func (s *PDS) Shutdown(ctx context.Context) error {
	_ = s.logger.Sync()

	s.Postgres.Close()
	s.logger.Info("postgres connections was closed")

	return nil
}

type GetPingRespData struct {
	Ok bool `json:"ok"`
}

func (s *PDS) Health() *GetPingRespData {
	return &GetPingRespData{Ok: true}
}

type GetInfoRespData struct {
	Name              string `json:"name"`
	PostgresIsRunning bool   `json:"postgres_is_running"`
}

func (s *PDS) Info(ctx context.Context) *GetInfoRespData {
	return &GetInfoRespData{
		Name:              string(ServiceName),
		PostgresIsRunning: s.Postgres.Ping(ctx) == nil,
	}
}
