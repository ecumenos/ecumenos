package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ecumenos/ecumenos/internal/toolkit/errorsutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/random"
	"github.com/ecumenos/ecumenos/internal/toolkit/timeutils"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/jackc/pgx/v4"
)

func (r *Repository) InsertOrbisSociusLaunchInvite(ctx context.Context, email, code string, osLaunchReID int64, expiredAt time.Time) (*models.OrbisSociusLaunchInvite, error) {
	id, err := random.GetSnowflakeID[models.OrbisSociusLaunchInvite](ctx, 0, r.GetOrbisSociusLaunchInviteByID)
	if err != nil {
		return nil, err
	}
	createdAt := time.Now()
	if expiredAt.Before(createdAt) {
		return nil, fmt.Errorf("expired at can not be before created at (expired at = %v, created at = %v)", timeutils.TimeToString(expiredAt), timeutils.TimeToString(createdAt))
	}

	query := fmt.Sprintf(`insert into public.orbes_socii_launch_invites
  (id, created_at, email, code, orbis_socius_launch_requests_id, expired_at)
  values ($1, $2, $3, $4, $5, $6);`)
	params := []interface{}{id, createdAt, email, code, osLaunchReID, expiredAt}
	if _, err := r.driver.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.OrbisSociusLaunchInvite{
		ID:                         id,
		CreatedAt:                  createdAt,
		Email:                      email,
		Code:                       code,
		OrbisSociusLaunchRequestID: osLaunchReID,
		ExpiredAt:                  expiredAt,
	}, nil
}

func scanRowOrbisSociusLaunchInvite(row pgx.Row) (*models.OrbisSociusLaunchInvite, error) {
	var i models.OrbisSociusLaunchInvite
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Email,
		&i.Code,
		&i.OrbisSociusLaunchRequestID,
		&i.ExpiredAt,
	)
	if err == nil {
		return &i, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) GetOrbisSociusLaunchInviteByID(ctx context.Context, id int64) (*models.OrbisSociusLaunchInvite, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, email, code, orbis_socius_launch_requests_id, expired_at
    from public.orbes_socii_launch_invites
		where id=$1;
	`)
	row, err := r.driver.QueryRow(ctx, q, id)
	if err != nil {
		return nil, err
	}

	return scanRowOrbisSociusLaunchInvite(row)
}

func (r *Repository) InsertOrbisSociusLaunchRequest(ctx context.Context, email, region, name, desc string, status models.OrbisSociusLaunchRequestStatus) (*models.OrbisSociusLaunchRequest, error) {
	id, err := random.GetSnowflakeID[models.OrbisSociusLaunchRequest](ctx, 0, r.GetOrbisSociusLaunchRequestByID)
	if err != nil {
		return nil, err
	}
	createdAt := time.Now()

	query := fmt.Sprintf(`insert into public.orbes_socii_launch_requests
  (id, created_at, email, region, orbis_socius_name, orbis_socius_description, status)
  values ($1, $2, $3, $4, $5, $6, $7);`)
	params := []interface{}{id, createdAt, email, region, name, desc, status}
	if _, err := r.driver.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.OrbisSociusLaunchRequest{
		ID:                     id,
		CreatedAt:              createdAt,
		Email:                  email,
		Region:                 region,
		OrbisSociusName:        name,
		OrbisSociusDescription: desc,
		Status:                 status,
	}, nil
}

func scanRowOrbisSociusLaunchRequest(row pgx.Row) (*models.OrbisSociusLaunchRequest, error) {
	var r models.OrbisSociusLaunchRequest
	err := row.Scan(
		&r.ID,
		&r.CreatedAt,
		&r.Email,
		&r.Region,
		&r.OrbisSociusName,
		&r.OrbisSociusDescription,
		&r.Status,
	)
	if err == nil {
		return &r, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) GetOrbisSociusLaunchRequestByID(ctx context.Context, id int64) (*models.OrbisSociusLaunchRequest, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, email, region, orbis_socius_name, orbis_socius_description, status
    from public.orbes_socii_launch_requests
		where id=$1;
	`)
	row, err := r.driver.QueryRow(ctx, q, id)
	if err != nil {
		return nil, err
	}

	return scanRowOrbisSociusLaunchRequest(row)
}

func (r *Repository) InsertOrbisSociusStats(ctx context.Context, orbisSociusID int64, alive bool) (*models.OrbisSociusStats, error) {
	id, err := random.GetSnowflakeID[models.OrbisSociusStats](ctx, 0, r.GetOrbisSociusStatsByID)
	if err != nil {
		return nil, err
	}
	createdAt := time.Now()

	query := fmt.Sprintf(`insert into public.orbes_socii_stats
  (id, created_at, orbis_socius_id, alive)
  values ($1, $2, $3, $4);`)
	params := []interface{}{id, createdAt, orbisSociusID, alive}
	if _, err := r.driver.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.OrbisSociusStats{
		ID:            id,
		CreatedAt:     createdAt,
		OrbisSociusID: orbisSociusID,
		Alive:         alive,
	}, nil
}

func scanRowOrbisSociusStats(row pgx.Row) (*models.OrbisSociusStats, error) {
	var s models.OrbisSociusStats
	err := row.Scan(
		&s.ID,
		&s.CreatedAt,
		&s.OrbisSociusID,
		&s.Alive,
	)
	if err == nil {
		return &s, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) GetOrbisSociusStatsByID(ctx context.Context, id int64) (*models.OrbisSociusStats, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, orbis_socius_id, alive
    from public.orbes_socii_stats
		where id=$1;
	`)
	row, err := r.driver.QueryRow(ctx, q, id)
	if err != nil {
		return nil, err
	}

	return scanRowOrbisSociusStats(row)
}

func (r *Repository) InsertOrbisSocius(ctx context.Context, ownerEmail, region, name, desc, url, apiKey string) (*models.OrbisSocius, error) {
	id, err := random.GetSnowflakeID[models.OrbisSocius](ctx, 0, r.GetOrbisSociusByID)
	if err != nil {
		return nil, err
	}
	createdAt := time.Now()
	updatedAt := time.Now()
	tombstoned := false

	query := fmt.Sprintf(`insert into public.orbes_socii
  (id, created_at, updated_at, tombstoned, owner_email, region, name, description, url, api_key)
  values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`)
	params := []interface{}{id, createdAt, updatedAt, tombstoned, ownerEmail, region, name, desc, url, apiKey}
	if _, err := r.driver.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.OrbisSocius{
		ID:          id,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Tombstoned:  tombstoned,
		OwnerEmail:  ownerEmail,
		Region:      region,
		Name:        name,
		Description: desc,
		URL:         url,
		APIKey:      apiKey,
	}, nil
}

func scanRowOrbisSocius(row pgx.Row) (*models.OrbisSocius, error) {
	var os models.OrbisSocius
	err := row.Scan(
		&os.ID,
		&os.CreatedAt,
		&os.UpdatedAt,
		&os.DeletedAt,
		&os.Tombstoned,
		&os.OwnerEmail,
		&os.Alive,
		&os.RobustnessStatus,
		&os.LastPingedAt,
		&os.Region,
		&os.Name,
		&os.Description,
		&os.URL,
		&os.APIKey,
	)
	if err == nil {
		return &os, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) GetOrbisSociusByID(ctx context.Context, id int64) (*models.OrbisSocius, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, deleted_at, tombstoned, owner_email, alive, robustness_status, last_pinged_at, region, name, description, url, api_key
    from public.orbes_socii
		where id=$1 and tombstoned=false;
	`)
	row, err := r.driver.QueryRow(ctx, q, id)
	if err != nil {
		return nil, err
	}

	return scanRowOrbisSocius(row)
}
