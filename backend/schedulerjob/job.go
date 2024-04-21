package schedulerjob

import (
	"gorm.io/gorm"
	"time"
	// "rockbackup/backend/host"
	// "rockbackup/backend/repository"
)

const (
	SchedulerJobStatusQueued    = "queued"
	SchedulerJobStatusRunning   = "running"
	SchedulerJobStatusCompleted = "completed"
)

func NewBackupJob(policyID uint, backupType, operator string) Job {
	queueTime := time.Now()
	priority := uint(5)

	return Job{
		QueueTime:  &queueTime,
		Status:     SchedulerJobStatusQueued,
		Priority:   priority,
		Operator:   operator,
		InSchedule: true,
	}
}

type Job struct {
	gorm.Model
	ID           uint       `gorm:"column:id;primaryKey;autoIncrement"`
	SourceHostID uint       `gorm:"column:source_host_id"`
	QueueTime    *time.Time `gorm:"column:queue_time"`
	StartTime    *time.Time `gorm:"column:start_time"`
	EndTime      *time.Time `gorm:"column:end_time"`
	Status       string     `gorm:"column:status"`
	Hostname     string     `gorm:"column:hostname"`
	RepositoryID uint       `gorm:"column:repository_id"`
	Priority     uint       `gorm:"column:priority"`
	InSchedule   bool       `gorm:"column:in_schedule"`
	Operator     string     `gorm:"column:operator"`
	// Repository   repository.Repository `gorm:"column:repository"`
	// SourceHost   SourceHost
	// Host         host.Host
}
