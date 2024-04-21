package job

import (
	"rockbackup/backend/agentd"
	"rockbackup/backend/scheduler"
	"rockbackup/backend/service"
)

type FileSchedulerJob struct {
	scheduler.SchedulerJob
}

func (j *FileSchedulerJob) Run(policy service.Policy, agent agentd.Agent) error {

	// agent is assigned

	// task1 agent is ocupied

	// task 1 is done agent is rleased

	// task 2

	// agent is released

	return nil
}
