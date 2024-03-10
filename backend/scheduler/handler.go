package scheduler

import (
	"github.com/hibiken/asynq"
)

func NewHandler(client *asynq.Client) *AsyncHandler {
	return &AsyncHandler{asynq: client}
}

type AsyncHandler struct {
	asynq *asynq.Client
}

func (h *AsyncHandler) StartBackup(policyID uint) error {
	return nil
}

func (h *AsyncHandler) StartRestore(policyID uint) error {
	return nil
}
