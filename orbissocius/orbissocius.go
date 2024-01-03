package orbissocius

import (
	"context"

	"github.com/ecumenos/fxecumenos"
	"github.com/ecumenos/fxecumenos/fxpostgres/postgres"
	"go.uber.org/zap"
)

var (
	ServiceName    fxecumenos.ServiceName = "orbis-socius"
	ServiceVersion fxecumenos.Version     = "v0.0.0"
)

type Config struct {
	Addr        string
	Prod        bool
	PostgresURL string
}

type OrbisSocius struct {
	Postgres *postgres.Driver
	logger   *zap.Logger
}

func New(cfg *Config, l *zap.Logger) (*OrbisSocius, error) {
	driver, err := postgres.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	return &OrbisSocius{
		Postgres: driver,
		logger:   l,
	}, nil
}

func (o *OrbisSocius) Start(ctx context.Context) error {
	if err := o.Postgres.Ping(ctx); err != nil {
		return err
	}
	o.logger.Info("postgres is started")

	return nil
}

func (o *OrbisSocius) Shutdown(ctx context.Context) error {
	_ = o.logger.Sync()

	o.Postgres.Close()
	o.logger.Info("postgres connections was closed")

	return nil
}

type GetPingRespData struct {
	Ok bool `json:"ok"`
}

func (o *OrbisSocius) Health() *GetPingRespData {
	return &GetPingRespData{Ok: true}
}

type GetInfoRespData struct {
	Name              string `json:"name"`
	PostgresIsRunning bool   `json:"postgres_is_running"`
}

func (o *OrbisSocius) Info(ctx context.Context) *GetInfoRespData {
	return &GetInfoRespData{
		Name:              string(ServiceName),
		PostgresIsRunning: o.Postgres.Ping(ctx) == nil,
	}
}
