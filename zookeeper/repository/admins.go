package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ecumenos/ecumenos/internal/toolkit/errorsutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/random"
	"github.com/ecumenos/ecumenos/models/common"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/jackc/pgx/v4"
)

func (r *Repository) InsertAdmin(ctx context.Context, email, passwordHash string) (*models.Admin, error) {
	id, err := random.GetSnowflakeID[models.Admin](ctx, 0, r.GetAdminByID)
	if err != nil {
		return nil, err
	}
	if !common.EmailRegex.MatchString(email) {
		return nil, fmt.Errorf("invalid email. it doesn't fulfill validation (email = %v)", email)
	}
	a, err := r.GetAdminByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if a != nil {
		return nil, fmt.Errorf("admin with this email has already exists (email = %v, ID = %v)", email, a.ID)
	}
	createdAt := time.Now()
	updatedAt := time.Now()
	tombstoned := false

	query := `insert into public.admins
  (id, created_at, updated_at, tombstoned, email, password_hash)
  values ($1, $2, $3, $4, $5, $6);`
	params := []interface{}{id, createdAt, updatedAt, tombstoned, email, passwordHash}
	if _, err := r.driver.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.Admin{
		ID:           id,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Tombstoned:   tombstoned,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}

func (r *Repository) GetAdminByID(ctx context.Context, id int64) (*models.Admin, error) {
	q := `
  select
    id, created_at, updated_at, deleted_at, tombstoned, email, password_hash
  from public.admins
  where id=$1 and tombstoned=false;`
	row, err := r.driver.QueryRow(ctx, q, id)
	if err != nil {
		return nil, err
	}

	var a models.Admin
	err = row.Scan(
		&a.ID,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
		&a.Tombstoned,
		&a.Email,
		&a.PasswordHash,
	)
	if err == nil {
		return &a, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error) {
	q := `
  select
    id, created_at, updated_at, deleted_at, tombstoned, email, password_hash
  from public.admins
  where email=$1 and tombstoned=false;`
	row, err := r.driver.QueryRow(ctx, q, email)
	if err != nil {
		return nil, err
	}

	var a models.Admin
	err = row.Scan(
		&a.ID,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
		&a.Tombstoned,
		&a.Email,
		&a.PasswordHash,
	)
	if err == nil {
		return &a, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) AssignRoleForAdmin(ctx context.Context, receiverAdminID, adminRoleID int64, granterID *int64) (*models.AdminsAdminRolesRelation, error) {
	granterAdminID := sql.NullInt64{}
	if granterID != nil {
		granterAdminID.Int64 = *granterID
		granterAdminID.Valid = true
	}
	grantedAt := time.Now()

	query := `insert into public.admins_admin_roles_relations
  (receiver_admin_id, granter_admin_id, role_id, granted_at)
  values ($1, $2, $3, $4);`
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
