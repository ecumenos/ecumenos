package zookeeper

import "time"

type OrbisSociusLaunchRequest struct {
	ID                     int64                          `json:"id"`
	CreatedAt              time.Time                      `json:"created_at"`
	ComptusID              int64                          `json:"comptus_id"`
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
