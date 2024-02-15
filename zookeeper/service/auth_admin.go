package service

import (
	"context"
	"fmt"

	"github.com/ecumenos/ecumenos/internal/toolkit/primitives"
)

func (s *Service) AuthorizeAdmin(ctx context.Context, token string) (int64, int64, error) {
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

func (s *Service) AuthorizeAdminWithRefreshToken(ctx context.Context, refreshToken string) (int64, int64, error) {
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
