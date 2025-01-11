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
	StarterDB
	repository.RepositoryDB
	schedulerjob.JobDB
}

type StarterDB interface {
	LoadJob(id uint) (*schedulerjob.Job, error)
	LoadRepository(id uint) (*repository.Repository, error)
	LoadPolicy(id uint) (*policy.Policy, error)
	LoadAgent(hostname string) (*agentd.Agent, error)
	// LoadBackend(id uint) (*repository.Backend, error)
}

type Starter struct {
	db DB
}

func (s *Starter) LoadPolicy(id uint) (*policy.Policy, error) {
	return s.db.LoadPolicy(id)
}
func (s *Starter) LoadAgent(hostname string) (*agentd.Agent, error) {
	return s.db.LoadAgent(hostname)
}

func (s *Starter) LoadJob(id uint) (*schedulerjob.Job, error) {
	return s.db.LoadJob(id)
}

func (f *Starter) LoadRepo(repoID uint) (*repository.Repository, error) {
	var repo *repository.Repository
	repo, err := f.db.LoadRepository(repoID)

	if err != nil {
		return nil, err
	}

	// repo = repository.NewRepository(name, backend, f.db)

	return repo, nil
}

func (s *Starter) StartFileBackupJobFile(ctx context.Context, p taskdef.BackupJobPayload) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("panic in StartBackupJobFile: %v\nStack Trace:\n%s", r, debug.Stack())
		}
	}()

	logger.Info("load job")
	job, err := s.db.LoadJob(p.ID)

	log := logger.WithFields(logrus.Fields{"job_id": p.ID})

	if err != nil {
		return err
	}

	log.Info("load policy")

	policy, err := s.LoadPolicy(job.PolicyID)
	if err != nil {
		return err
	}

	log.Info("load repo")
	repo, err := s.LoadRepo(policy.RepsoitoryID)
	if err != nil {
		return err
	}

	bset, err := repo.AddBackupset(s.db, job.ID, job.BackupType)

	if err != nil {
		return err
	}

	log.Info("load agent")
	agent, err := s.LoadAgent(job.Hostname)

	if err != nil {
		return err
	}

	log.Info("start to run job")
	filejob := fjob.NewFileBackupSchedulerJob(job, log)
	filejob.Run(ctx, s.db, policy, repo, agent, bset)

	return nil
}
