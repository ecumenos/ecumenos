package zookeeper

import (
	"database/sql"
	"time"
)

type OrbisSociusLaunchInvite struct {
	ID                         int64         `json:"id"`
	CreatedAt                  time.Time     `json:"created_at"`
	ComptusID                  int64         `json:"comptus_id"`
	AdminID                    int64         `json:"admin_id"`
	OrbisSociusID              sql.NullInt64 `json:"orbis_socius_id"`
	Code                       string        `json:"code"`
	APIKey                     string        `json:"api_key"`
	Used                       bool          `json:"used"`
	OrbisSociusLaunchRequestID sql.NullInt64 `json:"orbis_socius_launch_request_id"`
	ExpiredAt                  time.Time     `json:"expired_at"`
}
