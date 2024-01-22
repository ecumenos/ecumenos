package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ecumenos/ecumenos/internal/fxpostgres/postgres"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/ecumenos/ecumenos/zookeeper/config"
	"go.uber.org/zap"
)

type Repository struct {
	driver *postgres.Driver
	logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) (*Repository, error) {
	driver, err := postgres.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	return &Repository{
		driver: driver,
		logger: logger,
	}, nil
}

func (r *Repository) AssignRoleForAdmin(ctx context.Context, receiverAdminID, adminRoleID int64, granterID *int64) (*models.AdminsAdminRolesRelation, error) {
	granterAdminID := sql.NullInt64{}
	if granterID != nil {
		granterAdminID.Int64 = *granterID
		granterAdminID.Valid = true
	}
	grantedAt := time.Now()

	query := fmt.Sprintf(`insert into public.admins_admin_roles_relations
  (receiver_admin_id, granter_admin_id, role_id, granted_at)
  values ($1, $2, $3, $4);`)
	params := []interface{}{receiverAdminID, granterAdminID, adminRoleID, grantedAt}
	if _, err := r.driver.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.AdminsAdminRolesRelation{
		ReciverAdminID: receiverAdminID,
		GranterAdminID: granterAdminID,
		RoleID:         adminRoleID,
		GrantedAt:      grantedAt,
	}, nil
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.driver.Ping(ctx)
}
