package zookeeper

import (
	"database/sql"
	"regexp"
	"time"
)

type AdminRole struct {
	ID             int64                `json:"id"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	DeletedAt      sql.NullTime         `json:"deleted_at"`
	Tombstoned     bool                 `json:"tombstoned"`
	Name           string               `json:"name"`
	Permissions    AdminRolePermissions `json:"permissions"`
	CreatorAdminID int64                `json:"creator_admin_id"`
}

type AdminRolePermissions struct{}

var AdminRoleNameRegex = regexp.MustCompile("[a-z0-9_]+([a-z0-9_.-]+[a-z0-9_]+)?")
