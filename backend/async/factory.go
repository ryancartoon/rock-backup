package async

import (
	"context"
	"rockbackup/backend/agentd"
	"rockbackup/backend/async/taskdef"
	"rockbackup/backend/db"
	"rockbackup/backend/repository"
	"rockbackup/backend/schedulerjob"
	fjob "rockbackup/backend/schedulerjob/file"
	"rockbackup/backend/service"
)

var DB *db.DB

func initDB() {
	DB = db.InitTest()
}

type FactoryDB interface {
	LoadJob(id uint) (schedulerjob.Job, error)
	LoadRepository(id uint) (*repository.Repository, error)
	LoadPolicy(id uint) (service.Policy, error)
	LoadAgent(hostname string) (*agentd.Agent, error)
}

type Factory struct {
	db FactoryDB
}

func (f *Factory) StartBackupJobFile(ctx context.Context, p taskdef.BackupJobPayload, db schedulerjob.JobDB) error {
	job, err := f.db.LoadJob(p.ID)

	if err != nil {
		return err
	}

	policy, err := f.db.LoadPolicy(job.PolicyID)

	if err != nil {
		return err
	}

	repo, err := f.db.LoadRepository(policy.RepositoryID)

	if err != nil {
		return err
	}

	agent, err := f.db.LoadAgent(job.Hostname)

	if err != nil {
		return err
	}

	filejob := fjob.NewFileBackupSchedulerJob(job)
	filejob.Run(ctx, db, policy, repo, agent)

	return nil
}

func (f *Factory) LoadRepo(id uint) (*repository.Repository, error) {
	var repo *repository.Repository
	repo, err := f.db.LoadRepository(id)

	if err != nil {
		return nil, err
	}
	return repo, nil
}
