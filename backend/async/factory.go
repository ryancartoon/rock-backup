package async

import (
	"context"
	"rockbackup/backend/agentd"
	"rockbackup/backend/async/taskdef"
	"rockbackup/backend/policy"
	"rockbackup/backend/repository"
	"rockbackup/backend/schedulerjob"
	fjob "rockbackup/backend/schedulerjob/file"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

type DB interface {
	FactoryDB
	repository.RepositoryDB
}

type FactoryDB interface {
	LoadJob(id uint) (*schedulerjob.Job, error)
	LoadRepository(id uint) (*repository.Repository, error)
	LoadPolicy(id uint) (*policy.Policy, error)
	LoadAgent(hostname string) (*agentd.Agent, error)
}

type Factory struct {
	db DB
}

func (f *Factory) StartBackupJobFile(ctx context.Context, p taskdef.BackupJobPayload, db schedulerjob.JobDB) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("panic in StartBackupJobFile: %v\nStack Trace:\n%s", r, debug.Stack())
		}
	}()

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
	if err != nil {
		return err
	}

	repo, err := repository.LoadRepo(f.db, policy.RepositoryID)
	if err != nil {
		return err
	}

	bset, err := repo.AddBackupset(job.ID, job.BackupType)

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
	filejob.Run(ctx, db, policy, repo, agent, bset)

	return nil
}

func (f *Factory) LoadRepo(id uint) (*repository.Repository, error) {
	var repo *repository.Repository
	backend, err := f.db.LoadBackend(id)

	NewRepo(backend)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
