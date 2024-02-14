package zookeeper

import (
	"database/sql"
	"time"
)

type OrbisSocius struct {
	ID               int64            `json:"id"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	DeletedAt        sql.NullTime     `json:"deleted_at"`
	Tombstoned       bool             `json:"tombstoned"`
	OwnerComptusID   int64            `json:"owner_comptus_id"`
	ApproverAdminID  sql.NullInt64    `json:"approver_admin_id"`
	Alive            bool             `json:"alive"`
	RobustnessStatus RobustnessStatus `json:"robustness_status"`
	LastPingedAt     sql.NullTime     `json:"last_pinged_at"`
	Region           string           `json:"region"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	URL              string           `json:"url"`
	APIKey           string           `json:"api_key"`
}

type RobustnessStatus uint32

const (
	// At this stage, the system lacks basic robustness measures, making it susceptible to failures and vulnerabilities.
	Vulnerable RobustnessStatus = 0
	// The system has minimal error handling and fault tolerance, making it fragile in the face of disruptions.
	Fragile RobustnessStatus = 1
	// Basic robustness measures are in place, such as error handling and backup mechanisms, providing resilience against common failures.
	Resilient RobustnessStatus = 2
	// The system is designed with adaptability, allowing it to handle unexpected challenges and dynamically adjust to changes.
	Adaptable RobustnessStatus = 3
	// At this highest level, the system not only withstands disruptions but thrives on them, becoming antifragile by continuously improving and evolving in response to stressors.
	Antifragile RobustnessStatus = 4
)
