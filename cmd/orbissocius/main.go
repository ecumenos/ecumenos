package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ecumenos/ecumenos/internal/fxlogger"
	"github.com/ecumenos/ecumenos/internal/zerodowntime"
	"github.com/ecumenos/ecumenos/orbissocius"
	"github.com/ecumenos/ecumenos/orbissocius/admin"
	"github.com/ecumenos/ecumenos/orbissocius/app"
	"github.com/ecumenos/ecumenos/orbissocius/config"
	"go.uber.org/fx"

	cli "github.com/urfave/cli/v2"
)

func main() {
	if err := run(os.Args); err != nil {
		slog.Error("exiting", "err", err)
		os.Exit(-1)
	}
}

func run(args []string) error {
	app := cli.App{
		Name:    "api",
		Usage:   "serving API",
		Version: string(config.ServiceVersion),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "prod",
				Value:   false,
				EnvVars: []string{"PROD"},
			},
			&cli.StringFlag{
				Name:    "pg_url",
				Value:   "postgresql://ecumenosuser:rootpassword@localhost:5432/ecumenos_orbissociusdb",
				EnvVars: []string{"PG_URL"},
			},
		},
		Commands: []*cli.Command{
			runAppCmd,
			runAdminAppCmd,
			migrateUpCmd,
			migrateDownCmd,
		},
	}

	return app.Run(args)
}

type configuration struct {
	fx.Out

	Config       *config.Config
	LoggerConfig *fxlogger.Config
}

var runAppCmd = &cli.Command{
	Name:  "run-api-server",
	Usage: "run API HTTP server",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		return zerodowntime.HandleApp(fx.New(
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
			fx.Invoke(func(lc fx.Lifecycle, server *app.Server, shutdowner fx.Shutdowner) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							if err := server.Start(ctx); err != nil {
								slog.Error("zookeeper app server run error", "err", err)
								return
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return server.Shutdown(ctx)
					},
				})
			}),
		))
	},
}

var runAdminAppCmd = &cli.Command{
	Name:  "run-admin-server",
	Usage: "run Admin HTTP server",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		return zerodowntime.HandleApp(fx.New(
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
			fx.Invoke(func(lc fx.Lifecycle, adminServer *admin.Server, shutdowner fx.Shutdowner) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							if err := adminServer.Start(ctx); err != nil {
								slog.Error("zookeeper admin server run error", "err", err)
								return
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return adminServer.Shutdown(ctx)
					},
				})
			}),
		))
	},
}
