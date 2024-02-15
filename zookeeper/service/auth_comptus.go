package service

import (
	"context"
	"fmt"

	"github.com/ecumenos/ecumenos/internal/toolkit/primitives"
)

func (s *Service) AuthorizeComptus(ctx context.Context, token string) (int64, int64, error) {
	t, err := s.auth.DecodeToken(token)
	if err != nil {
		return 0, 0, err
	}
	subject := t.Subject()
	comptusID, err := primitives.StringToInt64(subject)
	if err != nil {
		return 0, 0, fmt.Errorf("token is corrupted (extracted subject = %v)", subject)
	}

	session, err := s.repo.GetComptusSessionByComptusIDAndToken(ctx, comptusID, token)
	if err != nil {
		return 0, 0, err
	}
	if session == nil {
		return 0, 0, fmt.Errorf("comptus session is not found (comptus id = %v)", comptusID)
	}

	return comptusID, session.ID, nil
}

func (s *Service) AuthorizeComptusWithRefreshToken(ctx context.Context, refreshToken string) (int64, int64, error) {
	t, err := s.auth.DecodeToken(refreshToken)
	if err != nil {
		return 0, 0, err
	}
	subject := t.Subject()
	comptusID, err := primitives.StringToInt64(subject)
	if err != nil {
		return 0, 0, fmt.Errorf("token is corrupted (extracted subject = %v)", subject)
	}

	session, err := s.repo.GetComptusSessionByComptusIDAndRefreshToken(ctx, comptusID, refreshToken)
	if err != nil {
		return 0, 0, err
	}
	if session == nil {
		return 0, 0, fmt.Errorf("comptus session is not found (comptus id = %v)", comptusID)
	}

	return comptusID, session.ID, nil
}
