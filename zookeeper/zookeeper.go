package zookeeper

import (
	"context"

	"github.com/ecumenos/fxecumenos"
	"github.com/ecumenos/fxecumenos/fxpostgres/postgres"
	"go.uber.org/zap"
)

var (
	ServiceName    fxecumenos.ServiceName = "zookeeper"
	ServiceVersion fxecumenos.Version     = "v0.0.0"
)

type Config struct {
	Addr        string
	Prod        bool
	PostgresURL string
}

type Zookeeper struct {
	Postgres *postgres.Driver
	logger   *zap.Logger
}

func New(cfg *Config, l *zap.Logger) (*Zookeeper, error) {
	driver, err := postgres.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	return &Zookeeper{
		Postgres: driver,
		logger:   l,
	}, nil
}

func (z *Zookeeper) Start(ctx context.Context) error {
	if err := z.Postgres.Ping(ctx); err != nil {
		return err
	}
	z.logger.Info("postgres is started")

	return nil
}

func (z *Zookeeper) Shutdown(ctx context.Context) error {
	_ = z.logger.Sync()

	z.Postgres.Close()
	z.logger.Info("postgres connections was closed")

	return nil
}

type GetPingRespData struct {
	Ok bool `json:"ok"`
}

func (z *Zookeeper) Health() *GetPingRespData {
	return &GetPingRespData{Ok: true}
}

type GetInfoRespData struct {
	Name              string `json:"name"`
	PostgresIsRunning bool   `json:"postgres_is_running"`
}

func (z *Zookeeper) Info(ctx context.Context) *GetInfoRespData {
	return &GetInfoRespData{
		Name:              string(ServiceName),
		PostgresIsRunning: z.Postgres.Ping(ctx) == nil,
	}
}
