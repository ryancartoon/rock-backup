package taskdef

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TaskTypeBackupJobFile = "rockbackup:backupjob-file"
)

type BackupJobPayload struct {
	ID uint
}

func NewBackupJobTask(id uint) (*asynq.Task, error) {
	payload, err := json.Marshal(BackupJobPayload{id})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskTypeBackupJobFile, payload), nil
}
