package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ecumenos/ecumenos/orbissocius"
	"github.com/ecumenos/fxecumenos/fxlogger/logger"
	"github.com/ecumenos/fxecumenos/zerodowntime"
	"go.uber.org/fx"
	"go.uber.org/zap"

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
		Version: string(orbissocius.ServiceVersion),
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
		Commands: []*cli.Command{runAppCmd, runAdminAppCmd},
	}

	return app.Run(args)
}

var runAppCmd = &cli.Command{
	Name:  "run-api-server",
	Usage: "run API HTTP server",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		return zerodowntime.HandleApp(fx.New(
			fx.Invoke(func(lc fx.Lifecycle, shutdowner fx.Shutdowner) {
				cfg := &orbissocius.Config{
					Prod:        cctx.Bool("prod"),
					Addr:        ":9091",
					PostgresURL: cctx.String("pg_url"),
				}
				l, err := newLogger(cctx.Bool("prod"))
				if err != nil {
					slog.Error("create logger error", "err", err)
					_ = shutdowner.Shutdown()
					return
				}
				instance, err := orbissocius.New(cfg, l)
				if err != nil {
					slog.Error("create orbis socius instance error", "err", err)
					_ = shutdowner.Shutdown()
					return
				}
				server := orbissocius.NewServer(cfg, instance, l)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							if err := instance.Start(ctx); err != nil {
								slog.Error("orbis socius instance run error", "err", err)
								return
							}
							if err := server.Start(ctx); err != nil {
								slog.Error("orbis socius server run error", "err", err)
								return
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						if err := server.Shutdown(ctx); err != nil {
							return err
						}

						return instance.Shutdown(ctx)
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
			fx.Invoke(func(lc fx.Lifecycle, shutdowner fx.Shutdowner) {
				cfg := &orbissocius.Config{
					Prod:        cctx.Bool("prod"),
					Addr:        ":9191",
					PostgresURL: cctx.String("pg_url"),
				}
				l, err := newLogger(cctx.Bool("prod"))
				if err != nil {
					slog.Error("create logger error", "err", err)
					_ = shutdowner.Shutdown()
					return
				}
				instance, err := orbissocius.New(cfg, l)
				if err != nil {
					slog.Error("create orbis socius instance error", "err", err)
					_ = shutdowner.Shutdown()
					return
				}
				server := orbissocius.NewAdminServer(cfg, instance, l)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							if err := instance.Start(ctx); err != nil {
								slog.Error("orbis socius instance run error", "err", err)
								return
							}
							if err := server.Start(ctx); err != nil {
								slog.Error("orbis socius admin server run error", "err", err)
								return
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						if err := server.Shutdown(ctx); err != nil {
							return err
						}

						return instance.Shutdown(ctx)
					},
				})
			}),
		))
	},
}

func newLogger(prod bool) (*zap.Logger, error) {
	var l *zap.Logger
	var err error
	if prod {
		l, err = logger.NewProductionLogger(string(orbissocius.ServiceName))
	} else {
		l, err = logger.NewDevelopmentLogger(string(orbissocius.ServiceName))
	}
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(l)

	return l, nil
}
