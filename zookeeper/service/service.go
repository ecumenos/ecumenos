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
	auth *Authorization
}

func New(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		auth: &Authorization{JWTSigningKey: cfg.JWTSecret},
	}
}

func (s *Service) CreateAdmin(ctx context.Context, email, password string) (*models.Admin, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return s.repo.InsertAdmin(ctx, email, passwordHash)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
		return errors.New("email is invalid")
	}
	if ok := checkPasswordHash(password, a.PasswordHash); !ok {
		return errors.New("password is invalid")
	}

	return nil
}

func (s *Service) GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error) {
	return s.repo.GetAdminByEmail(ctx, email)
}

func (s *Service) Authorize(ctx context.Context, token string) (int64, int64, error) {
	t, err := s.auth.DecodeToken(token)
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

func (s *Service) AuthorizeWithRefreshToken(ctx context.Context, refreshToken string) (int64, int64, error) {
	t, err := s.auth.DecodeToken(refreshToken)
	if err != nil {
		return 0, 0, err
	}
	subject := t.Subject()
	adminID, err := primitives.StringToInt64(subject)
	if err != nil {
		return 0, 0, fmt.Errorf("token is corrupted (extracted subject = %v)", subject)
	}

	session, err := s.repo.GetAdminSessionByAdminIDAndRefreshToken(ctx, adminID, refreshToken)
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
