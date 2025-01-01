package policy

import (
	"time"

	"rockbackup/backend/repository"

	"gorm.io/gorm"
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
	SourceType   string     `gorm:"column:source_type"                 json:"source_type"`
	SourceName   string     `gorm:"column:source_name"                 json:"source_name"`
	SourcePath   string     `gorm:"column:source_path"                 json:"source_path"`
	SourceHostID uint       `gorm:"column:source_host_id"              json:"source_host_id"`
	LastScanTime *time.Time `gorm:"column:last_scan_time"              json:"last_scan_time"`
}

// Policy backup policy
type Policy struct {
	gorm.Model
	ID             uint          `gorm:"column:id;primaryKey;autoIncrement"`
	Retention      uint          `gorm:"column:retention"`
	BackupSourceID uint          `gorm:"column:backup_source_id"`
	BackupSource   *BackupSource `gorm:"column:backup_source"`
	Hostname       string        `gorm:"column:hostname"`
	LogHostname    string        `gorm:"column:log_hostname"`
	Status         string        `gorm:"column:status"`
	RepositoryID   uint          `gorm:"column:repository_id"`
	ScheduleDesc   string        `gorm:"column:schedule_desc"`

	RepoName   string                 `gorm:"column:repo_name"`
	BackendID  uint                   `gorm:"column:backend_id"`
	Repository *repository.Repository `gorm:"-"`
	// BackupCycle  uint
}

type Instance struct {
	Name      string
	DataPath  string
	LoginPath string
	ConfPath  string
}
