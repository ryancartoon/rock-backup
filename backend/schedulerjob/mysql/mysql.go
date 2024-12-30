package mysql

import (
	"context"
	"rockbackup/backend/agentd"
	"rockbackup/backend/log"
	"rockbackup/backend/policy"
	"rockbackup/backend/repository"
	"rockbackup/backend/schedulerjob"
)

// job -> xtrabackup -> agent

var logger = log.New("agent.log")

func NewMysqlRestoreSchedulerJob(job schedulerjob.Job, t Tool, db DB, log *log.Logger) *MysqlBackupSchedulerJob {

	return &MysqlBackupSchedulerJob{job, t, log, db}
}

type Tool interface {
	Backup(
		ctx context.Context,
		agent *agentd.Agent,
		backupType string,
		instance policy.Instance,
		repo repository.Repository,
		targetPath string,
	) error
}

type DB interface {
}

type MysqlBackupSchedulerJob struct {
	schedulerjob.Job
	t      Tool
	logger *log.Logger
	db     DB
}

// func (j *MysqlBackupSchedulerJob) Run(
// 	ctx context.Context,
// 	db schedulerjob.JobDB,
// 	policy *policy.Policy,
// 	repo *repository.Repository,
// 	agent *agentd.Agent,
// 	bset *backupset.Backupset,
// ) error {

// 	if j.BackupType == "Full" {
// 		logger.Info("full backup, init repo")
// 	}

// 	size := j.t.Backup(ctx, agent, j.BackupType, instance, repo)
// }
