package zookeeper

import (
	"database/sql"
	"time"
)

type Comptus struct {
	ID           int64        `json:"id"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
	Tombstoned   bool         `json:"tombstoned"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"password_hash"`
	Patria       string       `json:"patria"`
	Lingua       string       `json:"lingua"`
}
