package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ecumenos/ecumenos/pds"
	"github.com/ecumenos/fxecumenos/zerodowntime"
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
		Version: string(pds.ServiceVersion),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "prod",
				Value:   false,
				EnvVars: []string{"PROD"},
			},
			&cli.StringFlag{
				Name:    "pg_url",
				Value:   "postgresql://ecumenosuser:rootpassword@localhost:5432/ecumenos_pdsdb",
				EnvVars: []string{"PG_URL"},
			},
		},
		Commands: []*cli.Command{runAppCmd},
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
				instance, err := pds.New(&pds.Config{
					Prod:        cctx.Bool("prod"),
					Addr:        ":9090",
					PostgresURL: cctx.String("pg_url"),
				})
				if err != nil {
					slog.Error("create PDS instance error", "err", err)
					_ = shutdowner.Shutdown()
				}
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							err := instance.Start(ctx)
							if err != nil {
								slog.Error("PDS run instance error", "err", err)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return instance.Shutdown(ctx)
					},
				})
			}),
		))
	},
}
