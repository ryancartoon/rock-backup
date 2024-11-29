package schedulerjob

import (
	"time"

	"gorm.io/gorm"
	// "rockbackup/backend/host"
	// "rockbackup/backend/repository"
)

const (
	SchedulerJobStatusCreated   = "created"
	SchedulerJobStatusQueued    = "queued"
	SchedulerJobStatusRunning   = "running"
	SchedulerJobStatusCompleted = "completed"

	JobTypeBackupFile = "backup_file"
)

type JobDB interface {
	SaveBackupResult(id uint, bsetID uint, snapID string, Size int64, FileNum int64) error
	SaveBackupError(id uint, err string)
}

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
	PolicyID     uint       `gorm:"column:policy_id"`
	SourceHostID uint       `gorm:"column:source_host_id"`
	QueueTime    *time.Time `gorm:"column:queue_time"`
	StartTime    *time.Time `gorm:"column:start_time"`
	EndTime      *time.Time `gorm:"column:end_time"`
	BackupType   string     `gorm:"column:backup_type"`
	JobType      string     `gorm:"column:job_type"`
	Status       string     `gorm:"column:status"`
	Hostname     string     `gorm:"column:hostname"`
	RepositoryID uint       `gorm:"column:repository_id"`
	Priority     uint       `gorm:"column:priority"`
	InSchedule   bool       `gorm:"column:in_schedule"`
	Operator     string     `gorm:"column:operator"`
	ErrorMessage string     `gorm:"column:error_message"`
	// ExternalBackupsetID string     `gorm:"column:external_backupset_id"`
	// Size                int64      `gorm:"column:size"`
	// FileNum             int64      `gorm:"column:file_num"`
	// Repository   repository.Repository `gorm:"column:repository"`
	// SourceHost   SourceHost
	// Host         host.Host
}
