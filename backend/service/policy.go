package service

import (
	"time"

	"gorm.io/gorm"
	"rockbackup/backend/repository"
)

const (
	ServiceStatusServing = "serving"

	BackupSourceTypeMySQL = "mysql"
	BackupSourceTypeFile  = "file"
)

// BackupSource source of backup
type BackupSource struct {
	gorm.Model
	ID           uint `gorm:"primaryKey;autoIncrement"`
	SourceType   string
	Name         string
	DataPath     string
	Host         string
	LastScanTime *time.Time
}

// Policy backup policy
type Policy struct {
	gorm.Model
	ID           uint
	Retention    uint
	SourceID     uint
	HostID       uint
	Status       string
	RepositoryID uint
	Repository   repository.Repository
	// ScheduleDesc string
	// BackupCycle  uint
}
