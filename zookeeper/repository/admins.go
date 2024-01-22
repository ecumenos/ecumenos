package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ecumenos/ecumenos/internal/toolkit/errorsutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/random"
	commonModels "github.com/ecumenos/ecumenos/models"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/jackc/pgx/v4"
)

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
