package zookeeper

import (
	"database/sql"
	"time"
)

type OrbisSociusStat struct {
	ID            int64         `json:"id"`
	CreatedAt     time.Time     `json:"created_at"`
	OrbisSociusID sql.NullInt64 `json:"orbis_socius_id"`
	Alive         bool          `json:"alive"`
}
