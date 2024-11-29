package backupset

import (
	"gorm.io/gorm"
	"time"
)

type Backupset struct {
	gorm.Model
	ID                  uint
	ExpiredAt           *time.Time `gorm:"column:expire_at"`
	BackupType          string     `gorm:"column:backup_type"`
	ExternalBackupsetID string     `gorm:"column:external_backupset_id"`
	Size                uint64     `gorm:"column:size"`
	FileNum             uint64     `gorm:"column:file_num"`
	rentention          int64      `gorm:"column:retention"`
	BackupTime          *time.Time `gorm:"column:backup_time"`
	RepositoryID        uint       `gorm:"column:repository_id"`
	JobID               uint       `gorm:"column:job_id"`
}
