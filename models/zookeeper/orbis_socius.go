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
	OwnerEmail       string           `json:"owner_email"`
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

type OrbisSociusStats struct {
	ID            int64     `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	OrbisSociusID int64     `json:"orbis_socius_id"`
	Alive         bool      `json:"alive"`
}

type OrbisSociusLaunchRequest struct {
	ID                     int64                          `json:"id"`
	CreatedAt              time.Time                      `json:"created_at"`
	Email                  string                         `json:"email"`
	Region                 string                         `json:"region"`
	OrbisSociusName        string                         `json:"orbis_socius_name"`
	OrbisSociusDescription string                         `json:"orbis_socius_description"`
	OrbisSociusURL         string                         `json:"orbis_socius_url"`
	Status                 OrbisSociusLaunchRequestStatus `json:"status"`
}

type OrbisSociusLaunchRequestStatus uint32

const (
	PendingOrbisSociusLaunchRequest  OrbisSociusLaunchRequestStatus = 0
	ViewedOrbisSociusLaunchRequest   OrbisSociusLaunchRequestStatus = 1
	ApprovedOrbisSociusLaunchRequest OrbisSociusLaunchRequestStatus = 2
	RejectedOrbisSociusLaunchRequest OrbisSociusLaunchRequestStatus = 3
)

type OrbisSociusLaunchInvite struct {
	ID                         int64     `json:"id"`
	CreatedAt                  time.Time `json:"created_at"`
	Email                      string    `json:"email"`
	Code                       string    `json:"code"`
	APIKey                     string    `json:"api_key"`
	Used                       bool      `json:"used"`
	OrbisSociusLaunchRequestID int64     `json:"orbis_socius_launch_request_id"`
	ExpiredAt                  time.Time `json:"expired_at"`
}
