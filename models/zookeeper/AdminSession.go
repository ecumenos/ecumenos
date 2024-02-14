package zookeeper

import (
	"database/sql"
	"time"
)

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
