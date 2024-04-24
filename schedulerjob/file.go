package job

import (
	"rockbackup/backend/agentd"
	"rockbackup/backend/repository"
	"rockbackup/backend/scheduler"
	"rockbackup/backend/schedulerjob"
	"rockbackup/backend/service"
)

type DB interface {
}

func NewFileBackupSchedulerJob(job scheduler.SchedulerJob, db DB) *FileBackupSchedulerJob {
	restic := Restic{}
	return &FileBackupSchedulerJob{SchedulerJob: job, Restic: restic}
}

type FileBackupSchedulerJob struct {
	scheduler.SchedulerJob
	Restic
}

func (j *FileBackupSchedulerJob) Run(policy service.Policy, repo repository.Repository, agent agentd.Agent) error {

	if j.IsFullBackup {
		j.Restic.InitRepo()
	}

	j.Restic.StartBackup(agent, repo.String())

	// agent is assigned

	// task1 agent is ocupied

	// task 1 is done agent is rleased

	// task 2

	// agent is released

	return nil
}

type Restic struct{}

func (r *Restic) InitRepo(agent agentd.Agent, repo repository.Repository) error {
	rsCmd := ["restic", "init" "--repo", repo.String()]
	out, err := agent.Run(rsCmd)
}

func (r *Restic) StartBackup(agent agentd.Agent, repo string) {

}

func RunCmdAgent(agent agentd.Agent, cmd string, env map[string]string) ([]byte, error) {
	var out []byte
	return out, nil
}
