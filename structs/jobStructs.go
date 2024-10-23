package structs

import "time"

const (
	JOB_PENDING  = "pending"
	JOB_PROGRESS = "in progress"
	JOB_FAILED   = "failed"
	JOB_SUCCESS  = "success"
)

type Job struct {
	ID        string    `json:"id" gorm:"primaryKey;type:char(36);not null"`
	JobName   string    `json:"job_name" gorm:"not null"`
	Payload   string    `json:"payload" gorm:"not null"`
	Status    string    `json:"status" gorm:"not null"`
	Attempts  int       `json:"attempts" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
