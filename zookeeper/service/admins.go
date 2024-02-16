package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ecumenos/ecumenos/models/common"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
)

func (s *Service) CreateAdmin(ctx context.Context, email, password string) (*models.Admin, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return s.repo.InsertAdmin(ctx, email, passwordHash)
}

func (s *Service) ValidateAdminCredentials(ctx context.Context, email, password string) error {
	if !common.EmailRegex.MatchString(email) {
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

func (s *Service) AssignRoleForAdmin(ctx context.Context, adminID, roleID int64, granterID *int64) error {
	_, err := s.repo.AssignRoleForAdmin(ctx, adminID, roleID, granterID)
	return err
}

func (s *Service) GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error) {
	return s.repo.GetAdminByEmail(ctx, email)
}
