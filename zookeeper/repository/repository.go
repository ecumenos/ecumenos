package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ecumenos/ecumenos/internal/fxpostgres/postgres"
	"github.com/ecumenos/ecumenos/internal/toolkit/errorsutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/random"
	"github.com/ecumenos/ecumenos/internal/toolkit/timeutils"
	commonModels "github.com/ecumenos/ecumenos/models"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/ecumenos/ecumenos/zookeeper/config"
	"github.com/jackc/pgx/v4"
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

func (r *Repository) InsertAdmin(ctx context.Context, email, passwordHash string) (*models.Admin, error) {
	id, err := random.GetSnowflakeID[models.Admin](ctx, 0, r.GetAdminByID)
	if err != nil {
		return nil, err
	}
	if !commonModels.EmailRegex.MatchString(email) {
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

	query := fmt.Sprintf(`insert into public.admins
  (id, created_at, updated_at, tombstoned, email, password_hash)
  values ($1, $2, $3, $4, $5, $6);`)
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
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, deleted_at, tombstoned, email, password_hash
    from public.admins
		where id=$1 and tombstoned=false;
	`)
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
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, deleted_at, tombstoned, email, password_hash
    from public.admins
		where email=$1 and tombstoned=false;
	`)
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

func (r *Repository) InsertAdminSession(ctx context.Context, adminID int64, t, rt string, expiredAt time.Time) (*models.AdminSession, error) {
	id, err := random.GetSnowflakeID[models.AdminSession](ctx, 0, r.GetAdminSessionByID)
	if err != nil {
		return nil, err
	}
	createdAt := time.Now()
	updatedAt := time.Now()
	tombstoned := false
	if expiredAt.Before(createdAt) {
		return nil, fmt.Errorf("expired at can not be before created at (expired at = %v, created at = %v)", timeutils.TimeToString(expiredAt), timeutils.TimeToString(createdAt))
	}

	query := fmt.Sprintf(`insert into public.admin_sessions
  (id, created_at, updated_at, expired_at, tombstoned, admin_id, token, refresh_token)
  values ($1, $2, $3, $4, $5, $6, $7, $8);`)
	params := []interface{}{id, createdAt, updatedAt, expiredAt, tombstoned, adminID, t, rt}
	if _, err := r.driver.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.AdminSession{
		ID:           id,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		ExpiredAt:    expiredAt,
		Tombstoned:   tombstoned,
		AdminID:      adminID,
		Token:        t,
		RefreshToken: rt,
	}, nil
}

func (r *Repository) GetAdminSessionByID(ctx context.Context, id int64) (*models.AdminSession, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, expired_at, deleted_at, tombstoned, admin_id, token, refresh_token
    from public.admin_sessions
		where id=$1 and tombstoned=false;
	`)
	row, err := r.driver.QueryRow(ctx, q, id)
	if err != nil {
		return nil, err
	}

	var s models.AdminSession
	err = row.Scan(
		&s.ID,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.ExpiredAt,
		&s.DeletedAt,
		&s.Tombstoned,
		&s.AdminID,
		&s.Token,
		&s.RefreshToken,
	)
	if err == nil {
		return &s, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) GetAdminSessionByAdminIDAndToken(ctx context.Context, adminID int64, token string) (*models.AdminSession, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, expired_at, deleted_at, tombstoned, admin_id, token, refresh_token
    from public.admin_sessions
		where admin_id=$1 and token=$2 and tombstoned=false;
	`)
	row, err := r.driver.QueryRow(ctx, q, adminID, token)
	if err != nil {
		return nil, err
	}

	var s models.AdminSession
	err = row.Scan(
		&s.ID,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.ExpiredAt,
		&s.DeletedAt,
		&s.Tombstoned,
		&s.AdminID,
		&s.Token,
		&s.RefreshToken,
	)
	if err == nil {
		return &s, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.driver.Ping(ctx)
}
