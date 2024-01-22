package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ecumenos/ecumenos/internal/fxlogger"
	"github.com/ecumenos/ecumenos/internal/zerodowntime"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/ecumenos/ecumenos/zookeeper"
	"github.com/ecumenos/ecumenos/zookeeper/config"
	"github.com/ecumenos/ecumenos/zookeeper/service"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var runSeedsCmd = &cli.Command{
	Name:  "run-seeds",
	Usage: "run seeds for the service",
	Flags: []cli.Flag{},
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
			fx.Invoke(func(lc fx.Lifecycle, shutdowner fx.Shutdowner, l *zap.Logger, s *service.Service) {
				defer func() { _ = shutdowner.Shutdown() }()

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				if err := seeds(ctx, l, s); err != nil {
					slog.Error("run seeds error", "err", err)
					return
				}
			}),
		))
	},
}

func seeds(ctx context.Context, l *zap.Logger, s *service.Service) error {
	var (
		superAdminEmail = "superadmin@mail.com"
		admins          = map[string]struct {
			password string
			roles    []string
		}{
			superAdminEmail: {
				password: "Password123!",
				roles:    []string{"superadmin"},
			},
		}
		adminsMap = map[string]*models.Admin{}
	)

	var (
		roles = map[string]struct {
			creator     string
			permissions models.AdminRolePermissions
		}{
			"superadmin": {
				creator:     superAdminEmail,
				permissions: models.AdminRolePermissions{},
			},
		}
		rolesMap = map[string]*models.AdminRole{}
	)

	for email, data := range admins {
		var err error
		admin, err := s.CreateAdmin(ctx, email, data.password)
		if err != nil {
			return err
		}
		adminsMap[email] = admin
	}

	for roleName, data := range roles {
		var creatorID int64
		if data.creator != "" {
			admin, ok := adminsMap[data.creator]
			if !ok {
				return fmt.Errorf("admin is not found (email = %v)", data.creator)
			}
			creatorID = admin.ID
		}
		role, err := s.CreateAdminRole(ctx, roleName, data.permissions, creatorID)
		if err != nil {
			return err
		}
		rolesMap[roleName] = role
	}

	superadmin, ok := adminsMap[superAdminEmail]
	if !ok {
		return fmt.Errorf("superadmin is not found (email = %v)", superAdminEmail)
	}
	for email, data := range admins {
		for _, roleName := range data.roles {
			role, ok := rolesMap[roleName]
			if !ok {
				return fmt.Errorf("role is not found (name = %v)", roleName)
			}
			admin, ok := adminsMap[email]
			if !ok {
				return fmt.Errorf("admin is not found (email = %v)", email)
			}
			if err := s.AssignRoleForAdmin(ctx, admin.ID, role.ID, &superadmin.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
