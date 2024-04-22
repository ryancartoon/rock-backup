package async

import (
	"rockbackup/backend/agentd"
	"rockbackup/backend/db"
	"rockbackup/backend/repository"
	"rockbackup/backend/schedulerjob"
	"rockbackup/backend/service"
)

var DB *db.DB

func initDB() {
	DB = db.InitTest()
}

type JobDB interface {
	LoadRepository(id uint) (*repository.Repository, error)
	LoadPolicy(id uint) (service.Policy, error)
	LoadJob(id uint) (*schedulerjob.Job, error)
}

type Factory struct {
	db DB
}

func (f *Factory) StartBackupFile(id, policyID uint) error {
	policy, err := f.db.LoadPolicy(policyID)

	if err != nil {
		return err
	}

	repo, err := f.LoadRepo(policy.RepositoryID)

	if err != nil {
		return err
	}

	job := schedulerjob.NewFileBackupSchedulerJob(id)
	job.Run(policy, repo, agent)

	return nil
}

func (f *Factory) LoadAgent() {}

func (f *Factory) LoadRepo(id uint) (*repository.Repository, error) {
	var repo *repository.Repository
	repo, err := f.db.LoadRepository(id)

	if err != nil {
		return nil, err
	}
	return repo, nil
}
