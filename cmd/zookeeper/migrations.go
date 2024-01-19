package main

import (
	"log/slog"

	"github.com/ecumenos/ecumenos/internal/fxpostgres/migrations"
	"github.com/ecumenos/ecumenos/internal/zerodowntime"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var migrateUpCmd = &cli.Command{
	Name:  "migrate-up",
	Usage: "run migrations up",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		app := fx.New(
			fx.Invoke(func(shutdowner fx.Shutdowner) error {
				l, err := newLogger(cctx.Bool("prod"))
				if err != nil {
					slog.Error("create logger error", "err", err)
					_ = shutdowner.Shutdown()
					return err
				}

				migrationsPath := "file://zookeeper/migrations"
				fn := migrations.NewMigrateUpFunc()
				if !cctx.Bool("prod") {
					l.Info("runnning migrate up",
						zap.String("db_url", cctx.String("pg_url")),
						zap.String("source_path", migrationsPath))
				}
				return fn(migrationsPath, cctx.String("pg_url")+"?sslmode=disable", l, shutdowner)
			}),
		)

		return zerodowntime.HandleApp(app)
	},
}

var migrateDownCmd = &cli.Command{
	Name:  "migrate-down",
	Usage: "run migrations down",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		app := fx.New(
			fx.Invoke(func(shutdowner fx.Shutdowner) error {
				l, err := newLogger(cctx.Bool("prod"))
				if err != nil {
					slog.Error("create logger error", "err", err)
					_ = shutdowner.Shutdown()
					return err
				}

				migrationsPath := "file://zookeeper/migrations"
				fn := migrations.NewMigrateDownFunc()
				if !cctx.Bool("prod") {
					l.Info("runnning migrate down",
						zap.String("db_url", cctx.String("pg_url")),
						zap.String("source_path", migrationsPath))
				}
				return fn(migrationsPath, cctx.String("pg_url")+"?sslmode=disable", l, shutdowner)
			}),
		)

		return zerodowntime.HandleApp(app)
	},
}
