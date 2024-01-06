package zookeeper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	commonModels "github.com/ecumenos/ecumenos/models"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/ecumenos/ecumenos/utils/snowflake"
	"github.com/ecumenos/fxecumenos"
	"github.com/ecumenos/fxecumenos/fxpostgres/postgres"
	"github.com/ecumenos/go-toolkit/errorsutils"
	"github.com/ecumenos/go-toolkit/primitives"
	"github.com/ecumenos/go-toolkit/timeutils"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

var (
	ServiceName    fxecumenos.ServiceName = "zookeeper"
	ServiceVersion fxecumenos.Version     = "v0.0.0"
)

type Config struct {
	Addr        string
	Prod        bool
	PostgresURL string
	JWTSecret   []byte
}

type Zookeeper struct {
	Postgres *postgres.Driver
	logger   *zap.Logger
	auth     *authorization
}

func New(cfg *Config, l *zap.Logger) (*Zookeeper, error) {
	driver, err := postgres.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	return &Zookeeper{
		Postgres: driver,
		logger:   l,
		auth:     &authorization{jwtSigningKey: cfg.JWTSecret},
	}, nil
}

func (z *Zookeeper) Start(ctx context.Context) error {
	if err := z.Postgres.Ping(ctx); err != nil {
		return err
	}
	z.logger.Info("postgres is started")

	return nil
}

func (z *Zookeeper) Shutdown(ctx context.Context) error {
	_ = z.logger.Sync()

	z.Postgres.Close()
	z.logger.Info("postgres connections was closed")

	return nil
}

type GetPingRespData struct {
	Ok bool `json:"ok"`
}

func (z *Zookeeper) Health() *GetPingRespData {
	return &GetPingRespData{Ok: true}
}

type GetInfoRespData struct {
	Name              string `json:"name"`
	PostgresIsRunning bool   `json:"postgres_is_running"`
}

func (z *Zookeeper) Info(ctx context.Context) *GetInfoRespData {
	return &GetInfoRespData{
		Name:              string(ServiceName),
		PostgresIsRunning: z.Postgres.Ping(ctx) == nil,
	}
}

func (z *Zookeeper) CreateAdmin(ctx context.Context, email, password string) (*models.Admin, error) {
	passwordHash, err := getPasswordHash(password)
	if err != nil {
		return nil, err
	}

	return z.insertAdmin(ctx, email, passwordHash)
}

var staticSalt = "aZedf4a"

func hash(in string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(staticSalt+in), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func getPasswordHash(password string) (string, error) {
	if err := models.VerifyPassword(password); err != nil {
		return "", err
	}
	return hash(password)
}

func validatePassword(password, passwordHash string) error {
	hashed, err := getPasswordHash(password)
	if err != nil {
		return err
	}

	if hashed != passwordHash {
		return errors.New("passwords doesn't match")
	}

	return nil
}

func (z *Zookeeper) CreateAdminRole(ctx context.Context, name string, permissions models.AdminRolePermissions, creatorID int64) (*models.AdminRole, error) {
	return z.insertAdminRole(ctx, name, permissions, creatorID)
}

func (z *Zookeeper) insertAdminRole(ctx context.Context, name string, permissions models.AdminRolePermissions, creatorID int64) (*models.AdminRole, error) {
	id, err := snowflake.GetSnowflakeID[models.AdminRole](ctx, 0, z.getAdminRoleByID)
	if err != nil {
		return nil, err
	}
	if !models.AdminRoleNameRegex.MatchString(name) {
		return nil, fmt.Errorf("invalid role name. it doesn't fulfill validation (name = %v)", name)
	}
	r, err := z.getAdminRoleByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if r != nil {
		return nil, fmt.Errorf("role with this name has already exists (name = %v, ID = %v)", name, r.ID)
	}
	createdAt := time.Now()
	updatedAt := time.Now()
	tombstoned := false

	query := fmt.Sprintf(`insert into public.admin_roles
  (id, created_at, updated_at, tombstoned, name, permissions, creator_admin_id)
  values ($1, $2, $3, $4, $5, $6, $7);`)
	params := []interface{}{id, createdAt, updatedAt, tombstoned, name, permissions, creatorID}
	if _, err := z.Postgres.QueryRow(ctx, query, params...); err != nil {
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

func (z *Zookeeper) getAdminRoleByID(ctx context.Context, id int64) (*models.AdminRole, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, deleted_at, tombstoned, name, permissions, creator_admin_id
    from public.admin_roles
		where id=$1 and tombstoned=false;
	`)
	row, err := z.Postgres.QueryRow(ctx, q, id)
	if err != nil {
		return nil, err
	}

	var r models.AdminRole
	err = row.Scan(
		&r.ID,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.DeletedAt,
		&r.Tombstoned,
		&r.Name,
		&r.Permissions,
		&r.CreatorAdminID,
	)
	if err == nil {
		return &r, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (z *Zookeeper) getAdminRoleByName(ctx context.Context, name string) (*models.AdminRole, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, deleted_at, tombstoned, name, permissions, creator_admin_id
    from public.admin_roles
		where name=$1 and tombstoned=false;
	`)
	row, err := z.Postgres.QueryRow(ctx, q, name)
	if err != nil {
		return nil, err
	}

	var r models.AdminRole
	err = row.Scan(
		&r.ID,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.DeletedAt,
		&r.Tombstoned,
		&r.Name,
		&r.Permissions,
		&r.CreatorAdminID,
	)
	if err == nil {
		return &r, nil
	}

	if errorsutils.Equals(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (z *Zookeeper) insertAdmin(ctx context.Context, email, passwordHash string) (*models.Admin, error) {
	id, err := snowflake.GetSnowflakeID[models.Admin](ctx, 0, z.getAdminByID)
	if err != nil {
		return nil, err
	}
	if !commonModels.EmailRegex.MatchString(email) {
		return nil, fmt.Errorf("invalid email. it doesn't fulfill validation (email = %v)", email)
	}
	a, err := z.getAdminByEmail(ctx, email)
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
	if _, err := z.Postgres.QueryRow(ctx, query, params...); err != nil {
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

func (z *Zookeeper) getAdminByID(ctx context.Context, id int64) (*models.Admin, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, deleted_at, tombstoned, email, password_hash
    from public.admins
		where id=$1 and tombstoned=false;
	`)
	row, err := z.Postgres.QueryRow(ctx, q, id)
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

func (z *Zookeeper) getAdminByEmail(ctx context.Context, email string) (*models.Admin, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, deleted_at, tombstoned, email, password_hash
    from public.admins
		where email=$1 and tombstoned=false;
	`)
	row, err := z.Postgres.QueryRow(ctx, q, email)
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

func (z *Zookeeper) AssignRoleForAdmin(ctx context.Context, adminID, roleID int64, granterID *int64) error {
	granterAdminID := sql.NullInt64{}
	if granterID != nil {
		granterAdminID.Int64 = *granterID
		granterAdminID.Valid = true
	}
	_, err := z.insertAdminsAdminRolesRelations(ctx, adminID, roleID, granterAdminID)
	return err
}

func (z *Zookeeper) insertAdminsAdminRolesRelations(ctx context.Context, receiverAdminID, adminRoleID int64, granterAdminID sql.NullInt64) (*models.AdminsAdminRolesRelation, error) {
	grantedAt := time.Now()

	query := fmt.Sprintf(`insert into public.admins_admin_roles_relations
  (receiver_admin_id, granter_admin_id, role_id, granted_at)
  values ($1, $2, $3, $4);`)
	params := []interface{}{receiverAdminID, granterAdminID, adminRoleID, grantedAt}
	if _, err := z.Postgres.QueryRow(ctx, query, params...); err != nil {
		return nil, err
	}

	return &models.AdminsAdminRolesRelation{
		ReciverAdminID: receiverAdminID,
		GranterAdminID: granterAdminID,
		RoleID:         adminRoleID,
		GrantedAt:      grantedAt,
	}, nil
}

func (z *Zookeeper) ValidateAdminCredentials(ctx context.Context, email, password string) error {
	if !commonModels.EmailRegex.MatchString(email) {
		return fmt.Errorf("invalid email. it doesn't fulfill validation (email = %v)", email)
	}
	a, err := z.getAdminByEmail(ctx, email)
	if err != nil {
		return err
	}
	if a == nil {
		return errors.New("email or password is invalid")
	}
	if err := validatePassword(password, a.PasswordHash); err != nil {
		return errors.New("email or password is invalid")
	}

	return nil
}

func (z *Zookeeper) CreateAdminSession(ctx context.Context, adminID int64) (*models.AdminSession, error) {
	a, err := z.getAdminByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, fmt.Errorf("failed create session for not existing admin (id = %v)", adminID)
	}

	tokExp, refTokExp := z.auth.getExpiredAt()
	token, refreshToken, err := z.auth.createAdminTokens(ctx, adminID, tokExp, refTokExp)
	if err != nil {
		return nil, err
	}

	return z.insertAdminSession(ctx, adminID, token, refreshToken, refTokExp)
}

func (z *Zookeeper) Authorize(ctx context.Context, token string) (int64, *models.AdminSession, error) {
	t, err := z.auth.decodeToken(token)
	if err != nil {
		return 0, nil, err
	}
	subject := t.Subject()
	adminID, err := primitives.StringToInt64(subject)
	if err != nil {
		return 0, nil, fmt.Errorf("token is corrupted (extracted subject = %v)", subject)
	}

	s, err := z.getAdminSessionByAdminIDAndToken(ctx, adminID, token)
	if err != nil {
		return 0, nil, err
	}

	return adminID, s, nil
}

func (z *Zookeeper) insertAdminSession(ctx context.Context, adminID int64, t, rt string, expiredAt time.Time) (*models.AdminSession, error) {
	id, err := snowflake.GetSnowflakeID[models.AdminSession](ctx, 0, z.getAdminSessionByID)
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
	if _, err := z.Postgres.QueryRow(ctx, query, params...); err != nil {
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

func (z *Zookeeper) getAdminSessionByID(ctx context.Context, id int64) (*models.AdminSession, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, expired_at, deleted_at, tombstoned, admin_id, token, refresh_token
    from public.admin_sessions
		where id=$1 and tombstoned=false;
	`)
	row, err := z.Postgres.QueryRow(ctx, q, id)
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

func (z *Zookeeper) getAdminSessionByAdminIDAndToken(ctx context.Context, adminID int64, token string) (*models.AdminSession, error) {
	q := fmt.Sprintf(`
		select
      id, created_at, updated_at, expired_at, deleted_at, tombstoned, admin_id, token, refresh_token
    from public.admin_sessions
		where admin_id=$1 and token=$2 and tombstoned=false;
	`)
	row, err := z.Postgres.QueryRow(ctx, q, adminID, token)
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
