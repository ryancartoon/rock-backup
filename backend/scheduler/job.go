package scheduler

import (
	"gorm.io/gorm"
	"time"

	"rockbackup/backend/host"
	"rockbackup/backend/repository"
)

const (
	SchedulerJobStatusQueued    = "queued"
	SchedulerJobStatusRunning   = "running"
	SchedulerJobStatusCompleted = "completed"
)

type SchedulerJob struct {
	gorm.Model
	ID           uint
	StartTime    *time.Time
	EndTime      *time.Time
	Status       string
	SourceHostID uint
	// SourceHost   SourceHost
	HostID       uint
	Host         host.Host
	RepositoryID uint
	Repository   repository.Repository
	Priority     uint
	InSchedule   bool
}
