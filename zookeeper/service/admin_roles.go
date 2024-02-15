package service

import (
	"context"

	models "github.com/ecumenos/ecumenos/models/zookeeper"
)

func (s *Service) CreateAdminRole(ctx context.Context, name string, permissions models.AdminRolePermissions, creatorID int64) (*models.AdminRole, error) {
	return s.repo.InsertAdminRole(ctx, name, permissions, creatorID)
}
