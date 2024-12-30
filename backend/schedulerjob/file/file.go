package filejob

import (
	"context"
	"rockbackup/backend/agentd"
	"rockbackup/backend/backupset"
	"rockbackup/backend/log"
	"rockbackup/backend/policy"
	"rockbackup/backend/repository"
	"rockbackup/backend/restic"
	"rockbackup/backend/schedulerjob"
)

var logger = log.New("agent.log")

func NewFileBackupSchedulerJob(job *schedulerjob.Job, log *log.Logger) *FileBackupSchedulerJob {
	envs := []string{"RESTIC_PASSWORD=redhat"}
	args := []string{"--json"}
	restic := restic.Restic{Name: "/usr/bin/restic", Envs: envs, GlobalArgs: args}

	return &FileBackupSchedulerJob{*job, restic, log}
}

type FileBackupSchedulerJob struct {
	schedulerjob.Job
	Restic restic.Restic
	logger *log.Logger
}

// func (j *FileBackupSchedulerJob) SaveBackupResult(id uint, bsetID uint, snapID string, Size int64, FileNum int64) error {
// 	return j.db.SaveBackupResult(id, bsetID, snapID, Size, FileNum)
// }
//
// func (j *FileBackupSchedulerJob ) SaveBackupError(id uint, err string) {
// 	return j.db.SaveBackupError(id, err
// }

func (j *FileBackupSchedulerJob) Run(
	ctx context.Context,
	db schedulerjob.JobDB,
	policy *policy.Policy,
	repo *repository.Repository,
	agent *agentd.Agent,
	bset *backupset.Backupset,
) error {
	var err error

	if j.BackupType == "Full" {
		logger.Info("full backup, init repo")
		err = j.Restic.InitRepo(ctx, agent, repo)

		if err != nil {
			return err
		}
	}

	logger.Info("start to run restic backup")
	snapID, size, fileNum, err := j.Restic.Backup(ctx, policy.BackupSource.SourcePath, agent, repo)

	logger.Infof("snap id is %s", snapID)
	logger.Infof("snap size is %d", size)

	if err != nil {
		db.SaveBackupError(j.ID, err.Error())
		return err
	}

	err = db.SaveBackupResult(j.ID, bset.ID, snapID, size, fileNum)

	if err != nil {
		return err
	}

	return nil
}

func NewFileRestoreSchedulerJob(job *schedulerjob.Job, log *log.Logger) *FileBackupSchedulerJob {
	envs := []string{"RESTIC_PASSWORD=redhat"}
	args := []string{"--json"}
	restic := restic.Restic{Name: "/usr/bin/restic", Envs: envs, GlobalArgs: args}

	return &FileBackupSchedulerJob{*job, restic, log}
}

type FileRestoreSchedulerJob struct {
	schedulerjob.Job
	Restic restic.Restic
	logger *log.Logger
}

func (j *FileRestoreSchedulerJob) Run(
	ctx context.Context,
	db schedulerjob.JobDB,
	repo *repository.Repository,
	agent *agentd.Agent,
	bset *backupset.Backupset,
	target string,
) error {

	logger.Info("start to run restic restore")
	err := j.Restic.Restore(ctx, agent, repo, bset, target)

	if err != nil {
		db.SaveBackupError(j.ID, err.Error())
		return err
	}

	return nil
}
