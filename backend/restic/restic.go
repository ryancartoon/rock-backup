package restic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"rockbackup/backend/agentd"
	"rockbackup/backend/backupset"
	"rockbackup/backend/repository"
)

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

func (r *Restic) InitRepo(ctx context.Context, agent *agentd.Agent, repo *repository.Repository) error {
	args := []string{"init", "--repo", repo.Backend.Path}
	rc, stdout, _, err := agent.RunCmd(ctx, r.Name, args, r.Envs)

	if err != nil {
		return err
	}

	if rc != 0 {
		return errors.New(string(stdout))
	}

	return nil
}

func (r *Restic) Backup(ctx context.Context, sourcePath string, agent *agentd.Agent, repo *repository.Repository) (string, int64, int64, error) {
	args := []string{"backup", sourcePath, "--repo", repo.GetTarget()}
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

func (r *Restic) Restore(ctx context.Context, agent *agentd.Agent, repo *repository.Repository, bset *backupset.Backupset, target string) error {
	args := []string{"restore", bset.ExternalBackupsetID, "--repo", repo.Backend.Path, "--target", target}
	args = append(args, r.GlobalArgs...)

	rc, stdout, _, err := agent.RunCmd(ctx, r.Name, args, r.Envs)

	if err != nil {
		return err
	}

	if rc != 0 {
		return errors.New(string(stdout))
	}

	return nil
}
