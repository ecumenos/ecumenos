package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ecumenos/ecumenos/internal/fxappsettings"
	"github.com/ecumenos/ecumenos/internal/fxlogger"
	"github.com/ecumenos/ecumenos/internal/zerodowntime"
	"github.com/ecumenos/ecumenos/zookeeper"
	"github.com/ecumenos/ecumenos/zookeeper/admin"
	"github.com/ecumenos/ecumenos/zookeeper/app"
	"github.com/ecumenos/ecumenos/zookeeper/config"
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
				Value:   "postgresql://ecumenosuser:rootpassword@localhost:5432/ecumenos_zookeeperdb",
				EnvVars: []string{"PG_URL"},
			},
		},
		Commands: []*cli.Command{
			runAppCmd,
			runAdminAppCmd,
			migrateUpCmd,
			migrateDownCmd,
			runSeedsCmd,
		},
	}

	return app.Run(args)
}

type configuration struct {
	fx.Out

	Config            *config.Config
	LoggerConfig      *fxlogger.Config
	AppSettingsConfig *fxappsettings.Config
}

var runAppCmd = &cli.Command{
	Name:  "run-api-server",
	Usage: "run API HTTP server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "jwt_secret",
			Usage:   "secret used for authenticating JWT tokens",
			Value:   "jwtsecretplaceholder",
			EnvVars: []string{"APP_JWT_SECRET"},
		},
		&cli.StringFlag{
			Name:    "locales_path",
			Usage:   "path to locales configuration",
			Value:   "./cmd/zookeeper/configurations/locales.yaml",
			EnvVars: []string{"APP_LOCALES_PATH"},
		},
		&cli.StringFlag{
			Name:    "regions_path",
			Usage:   "path to regions configuration",
			Value:   "./cmd/zookeeper/configurations/regions.yaml",
			EnvVars: []string{"APP_REGIONS_PATH"},
		},
	},
	Action: func(cctx *cli.Context) error {
		return zerodowntime.HandleApp(fx.New(
			fx.Options(fx.Provide(func() configuration {
				cfg := config.NewDefault()
				cfg.Prod = cctx.Bool("prod")
				cfg.PostgresURL = cctx.String("pg_url")
				cfg.JWTSecret = []byte(cctx.String("jwt_secret"))

				return configuration{
					Config:       cfg,
					LoggerConfig: &fxlogger.Config{Prod: cctx.Bool("prod")},
					AppSettingsConfig: &fxappsettings.Config{
						LocalesPath: cctx.String("locales_path"),
						RegionsPath: cctx.String("regions_path"),
					},
				}
			})),
			zookeeper.Module,
			fxlogger.Module,
			fxappsettings.Module,
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
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "jwt_secret",
			Usage:   "secret used for authenticating JWT tokens",
			Value:   "jwtsecretplaceholder",
			EnvVars: []string{"ADMIN_JWT_SECRET"},
		},
	},
	Action: func(cctx *cli.Context) error {
		return zerodowntime.HandleApp(fx.New(
			fx.Options(fx.Provide(func() configuration {
				cfg := config.NewDefault()
				cfg.Prod = cctx.Bool("prod")
				cfg.PostgresURL = cctx.String("pg_url")
				cfg.JWTSecret = []byte(cctx.String("jwt_secret"))

				return configuration{
					Config:       cfg,
					LoggerConfig: &fxlogger.Config{Prod: cctx.Bool("prod")},
				}
			})),
			zookeeper.Module,
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
