package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ecumenos/ecumenos/internal/toolkit/errorsutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/random"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/jackc/pgx/v4"
)

func (r *Repository) InsertAdminRole(ctx context.Context, name string, permissions models.AdminRolePermissions, creatorID int64) (*models.AdminRole, error) {
	id, err := random.GetSnowflakeID[models.AdminRole](ctx, 0, r.GetAdminRoleByID)
	if err != nil {
		return nil, err
	}
	if !models.AdminRoleNameRegex.MatchString(name) {
		return nil, fmt.Errorf("invalid role name. it doesn't fulfill validation (name = %v)", name)
	}
	role, err := r.GetAdminRoleByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if role != nil {
		return nil, fmt.Errorf("role with this name has already exists (name = %v, ID = %v)", name, role.ID)
	}
	createdAt := time.Now()
	updatedAt := time.Now()
	tombstoned := false

	query := fmt.Sprintf(`insert into public.admin_roles
  (id, created_at, updated_at, tombstoned, name, permissions, creator_admin_id)
  values ($1, $2, $3, $4, $5, $6, $7);`)
	params := []interface{}{id, createdAt, updatedAt, tombstoned, name, permissions, creatorID}
	if _, err := r.driver.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.AdminRole{
		ID:             id,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		Tombstoned:     tombstoned,
		Name:           name,
		Permissions:    permissions,
		CreatorAdminID: creatorID,
	}, nil
}

func (r *Repository) GetAdminRoleByID(ctx context.Context, id int64) (*models.AdminRole, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, deleted_at, tombstoned, name, permissions, creator_admin_id
    from public.admin_roles
		where id=$1 and tombstoned=false;
	`)
	row, err := r.driver.QueryRow(ctx, q, id)
	if err != nil {
		return nil, err
	}

	var role models.AdminRole
	err = row.Scan(
		&role.ID,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
		&role.Tombstoned,
		&role.Name,
		&role.Permissions,
		&role.CreatorAdminID,
	)
	if err == nil {
		return &role, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) GetAdminRoleByName(ctx context.Context, name string) (*models.AdminRole, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, deleted_at, tombstoned, name, permissions, creator_admin_id
    from public.admin_roles
		where name=$1 and tombstoned=false;
	`)
	row, err := r.driver.QueryRow(ctx, q, name)
	if err != nil {
		return nil, err
	}

	var role models.AdminRole
	err = row.Scan(
		&role.ID,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
		&role.Tombstoned,
		&role.Name,
		&role.Permissions,
		&role.CreatorAdminID,
	)
	if err == nil {
		return &role, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}
