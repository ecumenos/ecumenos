package service

import (
	"context"
	"errors"
	"fmt"

	commonModels "github.com/ecumenos/ecumenos/models"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
)

func (s *Service) CreateComptus(ctx context.Context, email, password, patria, lingua string) (*models.Comptus, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return s.repo.InsertComptus(ctx, email, passwordHash, patria, lingua)
}

func (s *Service) ValidateComptusCredentials(ctx context.Context, email, password string) error {
	if !commonModels.EmailRegex.MatchString(email) {
		return fmt.Errorf("invalid email. it doesn't fulfill validation (email = %v)", email)
	}
	a, err := s.repo.GetComptusByEmail(ctx, email)
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

func (s *Service) GetComptusByEmail(ctx context.Context, email string) (*models.Comptus, error) {
	return s.repo.GetComptusByEmail(ctx, email)
}

func (s *Service) GetComptusByID(ctx context.Context, id int64) (*models.Comptus, error) {
	return s.repo.GetComptusByID(ctx, id)
}
