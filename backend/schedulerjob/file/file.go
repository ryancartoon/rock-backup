package filejob

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"rockbackup/backend/agentd"
	"rockbackup/backend/log"
	"rockbackup/backend/policy"
	"rockbackup/backend/repository"
	"rockbackup/backend/schedulerjob"
)

var logger = log.New("agent.log")

var (
	ResticBackupFatalError  = errors.New("restic fatal error - no snapshot created")
	ResticBackupSourceError = errors.New("source data could not be read")
)

type Restic struct {
	Name       string
	Envs       []string
	GlobalArgs []string
}

// {"message_type":"summary","files_new":0,"files_changed":0,"files_unmodified":3,"dirs_new":0,"dirs_changed":0,
// "dirs_unmodified":3,"data_blobs":0,"tree_blobs":0,"data_added":0,"total_files_processed":3,
// "total_bytes_processed":11601,"total_duration":0.20655076,"snapshot_id":"6c2b23ec"}
type ResticBackupResponse struct {
	SnapshotID          string  `json:"snapshot_id"`
	DataAdded           int     `json:"data_added"`
	TotalFileProcessed  int     `json:"total_file_processed"`
	TotalBytesProcessed int     `json:"total_bytes_processed"`
	TotolDuration       float32 `json:"totol_duration"` // TODO check restic total duration data type
}

func NewFileBackupSchedulerJob(job *schedulerjob.Job, log *log.Logger) *FileBackupSchedulerJob {
	envs := []string{"RESTIC_PASSWORD=redhat"}
	args := []string{"--json"}
	restic := Restic{"/usr/bin/restic", envs, args}

	return &FileBackupSchedulerJob{*job, restic, log}
}

type FileBackupSchedulerJob struct {
	schedulerjob.Job
	Restic Restic
	logger *log.Logger
}

func (j *FileBackupSchedulerJob) Run(
	ctx context.Context,
	db schedulerjob.JobDB,
	policy *policy.Policy,
	repo *repository.Repository,
	agent *agentd.Agent,
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

	err = db.SaveBackupResult(j.ID, snapID, size, fileNum)

	if err != nil {
		return err
	}

	return nil
}

func (r *Restic) InitRepo(ctx context.Context, agent *agentd.Agent, repo *repository.Repository) error {
	args := []string{"init", "--repo", repo.MountPoint}
	rc, stdout, _, err := agent.RunCmd(ctx, r.Name, args, r.Envs)

	if err != nil {
		return err
	}

	if rc != 0 {
		return errors.New(string(stdout))
	}

	return nil
}

func (r *Restic) Backup(
	ctx context.Context,
	sourcePath string,
	agent *agentd.Agent,
	repo *repository.Repository,
) (string, int64, int64, error) {
	args := []string{"backup", sourcePath, "--repo", repo.MountPoint}
	args = append(args, r.GlobalArgs...)

	rc, stdout, _, err := agent.RunCmd(ctx, r.Name, args, r.Envs)

	if err != nil {
		return "", 0, 0, err
	}

	if rc == 1 {
		return "", 0, 0, ResticBackupFatalError
	}

	if rc == 3 {
		return "", 0, 0, ResticBackupSourceError
	}

	if rc != 0 {
		return "", 0, 0, errors.New(string(stdout))
	}

	lines := bytes.Split(stdout, []byte("\n"))
	summary := lines[len(lines)-2]

	resp := &ResticBackupResponse{}
	err = json.Unmarshal(summary, resp)

	if err != nil {
		return "", 0, 0, err
	}

	return resp.SnapshotID, int64(resp.TotalBytesProcessed), int64(resp.TotalFileProcessed), nil
}

func RunCmdAgent(agent agentd.Agent, cmd string, env map[string]string) ([]byte, error) {
	var out []byte
	return out, nil
}
