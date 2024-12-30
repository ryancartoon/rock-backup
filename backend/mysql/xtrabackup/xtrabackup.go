package xtrabackup

import (
	"context"
	"errors"
	"fmt"
	"rockbackup/backend/agentd"
	"rockbackup/backend/backupset"
	"rockbackup/backend/policy"
	"rockbackup/backend/repository"
)

type Xtrabackup struct {
	Version    string
	BinaryPath string
}

func (x *Xtrabackup) Backup(
	ctx context.Context,
	agent *agentd.Agent,
	backupType string,
	instance policy.Instance,
	repo repository.Repository,
	targetPath string,
) error {
	target := repo.MountPoint + targetPath
	var envs []string
	args := []string{
		"--backup",
		"--target-dir", target,
		"--login-path", instance.LoginPath,
		"--datadir", instance.DataPath,
		"--compress",
	}
	rc, stdout, _, err := agent.RunCmd(ctx, x.BinaryPath, args, envs)

	if err != nil {
		return err
	}

	x.ParseBackupOut(stdout)

	if rc != 0 {
		return errors.New(fmt.Sprintf("rc: %d", rc))
	}

	return nil
}

func (x *Xtrabackup) ParseBackupOut(out []byte) {}

func (x *Xtrabackup) Restore(
	ctx context.Context,
	agent *agentd.Agent,
	instance policy.Instance,
	repo repository.Repository,
	bset backupset.Backupset,
	targetPath string,

) error {
}


func (x *Xtrabackup) DeleteBackupset(bset backupset.Backupset) error {
	return nil
}
