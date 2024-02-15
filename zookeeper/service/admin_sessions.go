package service

import (
	"context"
	"fmt"

	models "github.com/ecumenos/ecumenos/models/zookeeper"
)

func (s *Service) CreateAdminSession(ctx context.Context, adminID int64) (*models.AdminSession, error) {
	a, err := s.repo.GetAdminByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, fmt.Errorf("failed create session for not existing admin (id = %v)", adminID)
	}

	tokExp, refTokExp := s.auth.GetExpiredAt()
	token, refreshToken, err := s.auth.CreateTokens(ctx, adminID, tokExp, refTokExp)
	if err != nil {
		return nil, err
	}

	return s.repo.InsertAdminSession(ctx, adminID, token, refreshToken, refTokExp)
}

func (s *Service) DeleteAdminSession(ctx context.Context, id int64) error {
	return s.repo.SetAdminSessionTombstonedByID(ctx, id)
}

func (s *Service) RefreshAdminSession(ctx context.Context, refreshToken string) (*models.AdminSession, error) {
	adminID, sessionID, err := s.AuthorizeAdminWithRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	tokExp, refTokExp := s.auth.GetExpiredAt()
	token, refreshToken, err := s.auth.CreateTokens(ctx, adminID, tokExp, refTokExp)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SetAdminSessionTokensByID(ctx, sessionID, token, refreshToken, refTokExp); err != nil {
		return nil, err
	}

	return s.repo.GetAdminSessionByID(ctx, sessionID)
}
