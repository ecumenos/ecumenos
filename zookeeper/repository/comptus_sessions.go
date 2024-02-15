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

func (r *Repository) InsertComptusSession(ctx context.Context, comptusID int64, t, rt string, expiredAt time.Time) (*models.ComptusSession, error) {
	id, err := random.GetSnowflakeID[models.ComptusSession](ctx, 0, r.GetComptusSessionByID)
	if err != nil {
		return nil, err
	}
	createdAt := time.Now()
	updatedAt := time.Now()
	tombstoned := false
	if expiredAt.Before(createdAt) {
		return nil, fmt.Errorf("expired at can not be before created at (expired at = %v, created at = %v)", timeutils.TimeToString(expiredAt), timeutils.TimeToString(createdAt))
	}

	query := fmt.Sprintf(`insert into public.comptus_sessions
  (id, created_at, updated_at, expired_at, tombstoned, comptus_id, token, refresh_token)
  values ($1, $2, $3, $4, $5, $6, $7, $8);`)
	params := []interface{}{id, createdAt, updatedAt, expiredAt, tombstoned, comptusID, t, rt}
	if _, err := r.driver.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.ComptusSession{
		ID:           id,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		ExpiredAt:    expiredAt,
		Tombstoned:   tombstoned,
		ComptusID:    comptusID,
		Token:        t,
		RefreshToken: rt,
	}, nil
}

func scanRowComptusSession(row pgx.Row) (*models.ComptusSession, error) {
	var s models.ComptusSession
	err := row.Scan(
		&s.ID,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.ExpiredAt,
		&s.DeletedAt,
		&s.Tombstoned,
		&s.ComptusID,
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

func (r *Repository) GetComptusSessionByID(ctx context.Context, id int64) (*models.ComptusSession, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, expired_at, deleted_at, tombstoned, comptus_id, token, refresh_token
    from public.comptus_sessions
		where id=$1 and tombstoned=false;
	`)
	row, err := r.driver.QueryRow(ctx, q, id)
	if err != nil {
		return nil, err
	}

	return scanRowComptusSession(row)
}

func (r *Repository) GetComptusSessionByComptusIDAndToken(ctx context.Context, comptusID int64, token string) (*models.ComptusSession, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, expired_at, deleted_at, tombstoned, comptus_id, token, refresh_token
    from public.comptus_sessions
		where comptus_id=$1 and token=$2 and tombstoned=false;
	`)
	row, err := r.driver.QueryRow(ctx, q, comptusID, token)
	if err != nil {
		return nil, err
	}

	return scanRowComptusSession(row)
}

func (r *Repository) GetComptusSessionByComptusIDAndRefreshToken(ctx context.Context, comptusID int64, refreshToken string) (*models.ComptusSession, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, expired_at, deleted_at, tombstoned, comptus_id, token, refresh_token
    from public.comptus_sessions
		where comptus_id=$1 and refresh_token=$2 and tombstoned=false;
	`)
	row, err := r.driver.QueryRow(ctx, q, comptusID, refreshToken)
	if err != nil {
		return nil, err
	}

	return scanRowComptusSession(row)
}

func (r *Repository) SetComptusSessionTombstonedByID(ctx context.Context, id int64) error {
	return r.driver.ExecuteQuery(ctx, "update public.comptus_sessions set tombstoned = true where id=$1", id)
}

func (r *Repository) SetComptusSessionTokensByID(ctx context.Context, id int64, t string, rt string, expiredAt time.Time) error {
	return r.driver.ExecuteQuery(ctx, "update public.comptus_sessions set updated_at = $2, expired_at = $3, token = $4, refresh_token = $5 where id=$1", id, time.Now(), expiredAt, t, rt)
}
