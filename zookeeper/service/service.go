package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ecumenos/ecumenos/internal/toolkit/primitives"
	commonModels "github.com/ecumenos/ecumenos/models"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/ecumenos/ecumenos/zookeeper/config"
	"github.com/ecumenos/ecumenos/zookeeper/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *repository.Repository
	auth *authorization
}

func New(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		auth: &authorization{jwtSigningKey: cfg.JWTSecret},
	}
}

func (s *Service) CreateAdmin(ctx context.Context, email, password string) (*models.Admin, error) {
	passwordHash, err := getPasswordHash(password)
	if err != nil {
		return nil, err
	}

	return s.repo.InsertAdmin(ctx, email, passwordHash)
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

func (s *Service) CreateAdminRole(ctx context.Context, name string, permissions models.AdminRolePermissions, creatorID int64) (*models.AdminRole, error) {
	return s.repo.InsertAdminRole(ctx, name, permissions, creatorID)
}

func (s *Service) AssignRoleForAdmin(ctx context.Context, adminID, roleID int64, granterID *int64) error {
	_, err := s.repo.AssignRoleForAdmin(ctx, adminID, roleID, granterID)
	return err
}

func (s *Service) ValidateAdminCredentials(ctx context.Context, email, password string) error {
	if !commonModels.EmailRegex.MatchString(email) {
		return fmt.Errorf("invalid email. it doesn't fulfill validation (email = %v)", email)
	}
	a, err := s.repo.GetAdminByEmail(ctx, email)
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

func (s *Service) CreateAdminSession(ctx context.Context, adminID int64) (*models.AdminSession, error) {
	a, err := s.repo.GetAdminByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, fmt.Errorf("failed create session for not existing admin (id = %v)", adminID)
	}

	tokExp, refTokExp := s.auth.getExpiredAt()
	token, refreshToken, err := s.auth.createAdminTokens(ctx, adminID, tokExp, refTokExp)
	if err != nil {
		return nil, err
	}

	return s.repo.InsertAdminSession(ctx, adminID, token, refreshToken, refTokExp)
}

func (s *Service) Authorize(ctx context.Context, token string) (int64, int64, error) {
	t, err := s.auth.decodeToken(token)
	if err != nil {
		return 0, 0, err
	}
	subject := t.Subject()
	adminID, err := primitives.StringToInt64(subject)
	if err != nil {
		return 0, 0, fmt.Errorf("token is corrupted (extracted subject = %v)", subject)
	}

	session, err := s.repo.GetAdminSessionByAdminIDAndToken(ctx, adminID, token)
	if err != nil {
		return 0, 0, err
	}
	if session == nil {
		return 0, 0, fmt.Errorf("admin session is not found (admin id = %v)", adminID)
	}

	return adminID, session.ID, nil
}

type PingServicesResult struct {
	PostgresIsRunning bool
}

func (s *Service) PingServices(ctx context.Context) *PingServicesResult {
	return &PingServicesResult{
		PostgresIsRunning: s.repo.Ping(ctx) == nil,
	}
}
