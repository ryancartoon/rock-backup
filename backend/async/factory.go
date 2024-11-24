package async

import (
	"context"
	"github.com/sirupsen/logrus"
	"rockbackup/backend/agentd"
	"rockbackup/backend/async/taskdef"
	"rockbackup/backend/policy"
	"rockbackup/backend/repository"
	"rockbackup/backend/schedulerjob"
	fjob "rockbackup/backend/schedulerjob/file"
)

type FactoryDB interface {
	LoadJob(id uint) (*schedulerjob.Job, error)
	LoadRepository(id uint) (*repository.Repository, error)
	LoadPolicy(id uint) (*policy.Policy, error)
	LoadAgent(hostname string) (*agentd.Agent, error)
}

type Factory struct {
	db FactoryDB
}

func (f *Factory) StartBackupJobFile(ctx context.Context, p taskdef.BackupJobPayload, db schedulerjob.JobDB) error {
	logger.Info("load job")
	job, err := f.db.LoadJob(p.ID)

	log := logger.WithFields(logrus.Fields{"job_id": p.ID})

	if err != nil {
		return err
	}

	log.Info("load policy")
	policy, err := f.db.LoadPolicy(job.PolicyID)

	if err != nil {
		return err
	}

	log.Info("load repo")
	policy, err = f.db.LoadPolicy(job.PolicyID)
	repo, err := f.db.LoadRepository(policy.RepositoryID)

	if err != nil {
		return err
	}

	log.Info("load agent")
	agent, err := f.db.LoadAgent(job.Hostname)

	if err != nil {
		return err
	}

	log.Info("start to run job")
	filejob := fjob.NewFileBackupSchedulerJob(job, log)
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
