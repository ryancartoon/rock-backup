package job

import (
	"rockbackup/backend/scheduler"
	"rockbackup/backend/service"
)

type MysqlBackupJob struct {
	scheduler.SchedulerJob
}

func (j MysqlBackupJob) Run(policy service.Policy, BackupStartType, BackupType string) {

	// LoginPath
	// DataPath
	// Retention
	// Version
	// ConfPath

}
