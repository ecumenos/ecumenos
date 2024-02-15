package service

import (
	"context"
	"fmt"

	models "github.com/ecumenos/ecumenos/models/zookeeper"
)

func (s *Service) CreateComptusSession(ctx context.Context, comptusID int64) (*models.ComptusSession, error) {
	c, err := s.repo.GetComptusByID(ctx, comptusID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("failed create session for not existing comptus (id = %v)", comptusID)
	}

	tokExp, refTokExp := s.auth.GetExpiredAt()
	token, refreshToken, err := s.auth.CreateTokens(ctx, comptusID, tokExp, refTokExp)
	if err != nil {
		return nil, err
	}

	return s.repo.InsertComptusSession(ctx, comptusID, token, refreshToken, refTokExp)
}

func (s *Service) DeleteComptusSession(ctx context.Context, id int64) error {
	return s.repo.SetComptusSessionTombstonedByID(ctx, id)
}

func (s *Service) RefreshComptusSession(ctx context.Context, refreshToken string) (*models.ComptusSession, error) {
	comptusID, sessionID, err := s.AuthorizeComptusWithRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	tokExp, refTokExp := s.auth.GetExpiredAt()
	token, refreshToken, err := s.auth.CreateTokens(ctx, comptusID, tokExp, refTokExp)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SetAdminSessionTokensByID(ctx, sessionID, token, refreshToken, refTokExp); err != nil {
		return nil, err
	}

	return s.repo.GetComptusSessionByID(ctx, sessionID)
}
