package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ecumenos/ecumenos/internal/zerodowntime"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/ecumenos/ecumenos/zookeeper"
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
			fx.Invoke(func(lc fx.Lifecycle, shutdowner fx.Shutdowner) {
				defer func() { _ = shutdowner.Shutdown() }()
				cfg := &zookeeper.Config{
					Prod:        cctx.Bool("prod"),
					PostgresURL: cctx.String("pg_url"),
				}
				l, err := newLogger(cctx.Bool("prod"))
				if err != nil {
					slog.Error("create logger error", "err", err)
					return
				}
				instance, err := zookeeper.New(cfg, l)
				if err != nil {
					slog.Error("create zookeeper instance error", "err", err)
					return
				}

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if err := instance.Start(ctx); err != nil {
					slog.Error("zookeeper instance run error", "err", err)
					return
				}

				if err := seeds(ctx, l, instance); err != nil {
					slog.Error("run seeds error", "err", err)
					return
				}
				if err := instance.Shutdown(ctx); err != nil {
					slog.Error("zookeeper instance run error", "err", err)
					return
				}
			}),
		))
	},
}

func seeds(ctx context.Context, l *zap.Logger, z *zookeeper.Zookeeper) error {
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
		admin, err := z.CreateAdmin(ctx, email, data.password)
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
		role, err := z.CreateAdminRole(ctx, roleName, data.permissions, creatorID)
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
			if err := z.AssignRoleForAdmin(ctx, admin.ID, role.ID, &superadmin.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
