package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ecumenos/ecumenos/internal/toolkit/errorsutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/random"
	"github.com/ecumenos/ecumenos/models/common"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/jackc/pgx/v4"
)

func (r *Repository) InsertComptus(ctx context.Context, email, passwordHash, patria, lingua string) (*models.Comptus, error) {
	id, err := random.GetSnowflakeID[models.Comptus](ctx, 0, r.GetComptusByID)
	if err != nil {
		return nil, err
	}
	if !common.EmailRegex.MatchString(email) {
		return nil, fmt.Errorf("invalid email. it doesn't fulfill validation (email = %v)", email)
	}
	c, err := r.GetAdminByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if c != nil {
		return nil, fmt.Errorf("comptus with this email has already exists (email = %v, ID = %v)", email, c.ID)
	}
	createdAt := time.Now()
	updatedAt := time.Now()
	tombstoned := false

	query := `insert into public.compti
  (id, created_at, updated_at, tombstoned, email, password_hash, patria, lingua)
  values ($1, $2, $3, $4, $5, $6, $7, $8);`
	params := []interface{}{id, createdAt, updatedAt, tombstoned, email, passwordHash, patria, lingua}
	if _, err := r.driver.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.Comptus{
		ID:           id,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Tombstoned:   tombstoned,
		Email:        email,
		PasswordHash: passwordHash,
		Patria:       patria,
		Lingua:       lingua,
	}, nil
}

func (r *Repository) GetComptusByID(ctx context.Context, id int64) (*models.Comptus, error) {
	q := `
  select
    id, created_at, updated_at, deleted_at, tombstoned, email, password_hash, patria, lingua
  from public.compti
  where id=$1 and tombstoned=false;`
	row, err := r.driver.QueryRow(ctx, q, id)
	if err != nil {
		return nil, err
	}

	var c models.Comptus
	err = row.Scan(
		&c.ID,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
		&c.Tombstoned,
		&c.Email,
		&c.PasswordHash,
		&c.Patria,
		&c.Lingua,
	)
	if err == nil {
		return &c, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) GetComptusByEmail(ctx context.Context, email string) (*models.Comptus, error) {
	q := `
  select
    id, created_at, updated_at, deleted_at, tombstoned, email, password_hash, patria, lingua
  from public.compti
  where email=$1 and tombstoned=false;`
	row, err := r.driver.QueryRow(ctx, q, email)
	if err != nil {
		return nil, err
	}

	var c models.Comptus
	err = row.Scan(
		&c.ID,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
		&c.Tombstoned,
		&c.Email,
		&c.PasswordHash,
		&c.Patria,
		&c.Lingua,
	)
	if err == nil {
		return &c, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}
