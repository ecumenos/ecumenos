package service

import "context"

func (s *Service) PingServices(ctx context.Context) *map[string]interface{} {
	return &map[string]interface{}{
		"postgres": s.repo.Ping(ctx) == nil,
	}
}
