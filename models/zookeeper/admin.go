package zookeeper

import (
	"database/sql"
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

type Admin struct {
	ID           int64        `json:"id"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
	Tombstoned   bool         `json:"tombstoned"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"password_hash"`
}

type AdminsAdminRolesRelation struct {
	ReciverAdminID int64         `json:"receiver_admin_id"`
	GranterAdminID sql.NullInt64 `json:"granter_admin_id"`
	RoleID         int64         `json:"role_id"`
	GrantedAt      time.Time     `json:"granted_at"`
}

type AdminSession struct {
	ID           int64        `json:"id"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	ExpiredAt    time.Time    `json:"expired_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
	Tombstoned   bool         `json:"tombstoned"`
	AdminID      int64        `json:"admin_id"`
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
}
