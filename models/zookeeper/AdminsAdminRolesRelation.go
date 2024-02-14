package zookeeper

import (
	"database/sql"
	"time"
)

type AdminsAdminRolesRelation struct {
	ReciverAdminID int64         `json:"receiver_admin_id"`
	GranterAdminID sql.NullInt64 `json:"granter_admin_id"`
	RoleID         int64         `json:"role_id"`
	GrantedAt      time.Time     `json:"granted_at"`
}
