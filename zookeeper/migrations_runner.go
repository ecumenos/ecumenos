package zookeeper

import (
	"github.com/ecumenos/ecumenos/internal/fxpostgres/migrations"
	"github.com/ecumenos/ecumenos/zookeeper/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MigrationsRunner struct {
	prod                   bool
	postgresURL            string
	postgresMigrationsPath string
	logger                 *zap.Logger
	shutdowner             fx.Shutdowner
}

func NewMigrationsRunner(cfg *config.Config, logger *zap.Logger, shutdowner fx.Shutdowner) *MigrationsRunner {
	return &MigrationsRunner{
		postgresURL:            cfg.PostgresURL,
		postgresMigrationsPath: cfg.PostgresMigrationsPath,
		logger:                 logger,
		shutdowner:             shutdowner,
	}
}

func (r *MigrationsRunner) MigrateUp() error {
	fn := migrations.NewMigrateUpFunc()
	if !r.prod {
		r.logger.Info("runnning migrate up",
			zap.String("db_url", r.postgresURL),
			zap.String("source_path", r.postgresMigrationsPath))
	}
	return fn(r.postgresMigrationsPath, r.postgresURL+"?sslmode=disable", r.logger, r.shutdowner)
}

func (r *MigrationsRunner) MigrateDown() error {
	fn := migrations.NewMigrateDownFunc()
	if !r.prod {
		r.logger.Info("runnning migrate down",
			zap.String("db_url", r.postgresURL),
			zap.String("source_path", r.postgresMigrationsPath))
	}
	return fn(r.postgresMigrationsPath, r.postgresURL+"?sslmode=disable", r.logger, r.shutdowner)
}
