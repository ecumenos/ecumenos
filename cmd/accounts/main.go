package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ecumenos/ecumenos/accounts"
	"github.com/ecumenos/ecumenos/accounts/app"
	"github.com/ecumenos/ecumenos/accounts/config"
	"github.com/ecumenos/ecumenos/internal/fxlogger"
	"github.com/ecumenos/ecumenos/internal/zerodowntime"
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
		},
		Commands: []*cli.Command{
			runAppCmd,
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

				return configuration{
					Config:       cfg,
					LoggerConfig: &fxlogger.Config{Prod: cctx.Bool("prod")},
				}
			})),
			accounts.Module,
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
