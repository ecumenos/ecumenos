package main

import (
	"github.com/ecumenos/ecumenos/internal/fxlogger"
	"github.com/ecumenos/ecumenos/internal/zerodowntime"
	"github.com/ecumenos/ecumenos/orbissocius"
	"github.com/ecumenos/ecumenos/orbissocius/config"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

var migrateUpCmd = &cli.Command{
	Name:  "migrate-up",
	Usage: "run migrations up",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		app := fx.New(
			fx.Options(fx.Provide(func() configuration {
				cfg := config.NewDefault()
				cfg.Prod = cctx.Bool("prod")
				cfg.PostgresURL = cctx.String("pg_url")

				return configuration{
					Config:       cfg,
					LoggerConfig: &fxlogger.Config{Prod: cctx.Bool("prod")},
				}
			})),
			orbissocius.Module,
			fxlogger.Module,
			fx.Invoke(func(runner *orbissocius.MigrationsRunner) error {
				return runner.MigrateUp()
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
			fx.Options(fx.Provide(func() configuration {
				cfg := config.NewDefault()
				cfg.Prod = cctx.Bool("prod")
				cfg.PostgresURL = cctx.String("pg_url")

				return configuration{
					Config:       cfg,
					LoggerConfig: &fxlogger.Config{Prod: cctx.Bool("prod")},
				}
			})),
			orbissocius.Module,
			fxlogger.Module,
			fx.Invoke(func(runner *orbissocius.MigrationsRunner) error {
				return runner.MigrateDown()
			}),
		)

		return zerodowntime.HandleApp(app)
	},
}
