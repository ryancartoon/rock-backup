package service

import (
	"time"

	"gorm.io/gorm"
	"rockbackup/backend/repository"
	// "rockbackup/backend/schedules"
)

const (
	ServiceStatusServing = "serving"

	BackupSourceTypeMySQL = "mysql"
	BackupSourceTypeFile  = "file"
)

// BackupSource source of backup
type BackupSource struct {
	gorm.Model
	ID           uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SourceType   string     `json:"source_type" gorm:"column:source_type"`
	SourceName   string     `json:"source_name" gorm:"column:source_name"`
	SourcePath   string     `json:"source_path" gorm:"column:source_path"`
	SourceHostID uint       `json:"source_host_id" gorm:"column:source_host_id"`
	LastScanTime *time.Time `json:"last_scan_time" gorm:"column:last_scan_time"`
}

// Policy backup policy
type Policy struct {
	gorm.Model
	ID             uint                   `gorm:"column:id;primaryKey;autoIncrement"`
	Retention      uint                   `gorm:"column:retention"`
	BackupSourceID uint                   `gorm:"column:backup_source_id"`
	BackupSource   *BackupSource          `gorm:"column:backup_source"`
	Hostname       string                 `gorm:"column:hostname"`
	Status         string                 `gorm:"column:status"`
	RepositoryID   uint                   `gorm:"column:repository_id"`
	Repository     *repository.Repository `gorm:"column:repository"`
	ScheduleDesc   string                 `gorm:"column:schedule_desc"`
	// BackupCycle  uint
}
