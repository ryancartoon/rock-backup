package scheduler

import (
	"rockbackup/backend/async/taskdef"
	"rockbackup/backend/schedulerjob"

	"github.com/hibiken/asynq"
)

func NewHandler(client *asynq.Client) *AsyncHandler {
	return &AsyncHandler{asynq: client}
}

type AsyncHandler struct {
	asynq *asynq.Client
}

func (h *AsyncHandler) Start(job JobInSchedule) error {

	if job.JobType == schedulerjob.JobTypeBackupFile {
		return h.StartBackup(job)
	}

	return nil
}

func (h *AsyncHandler) StartBackup(job JobInSchedule) error {
	t, err := taskdef.NewBackupJobTask(job.ID)

	if err != nil {
		return err
	}

	taskInfo, err := h.asynq.Enqueue(t)

	if err != nil {
		logger.Errorf("%v", taskInfo)
		return err
	}

	logger.Infof("%v", taskInfo)

	return nil
}

func (h *AsyncHandler) StartRestore(policyID uint) error {
	return nil
}
