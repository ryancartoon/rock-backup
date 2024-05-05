package scheduler

import (
	"errors"
	"fmt"
	"rockbackup/backend/backupset"
	"rockbackup/backend/repository"
	"rockbackup/backend/schedulerjob"
	"rockbackup/backend/service"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type JobInSchedule struct {
	schedulerjob.Job
	Repository repository.Repository
}

type DB interface {
	AddSchedulerJob(*schedulerjob.Job) error
	GetPolicy(uint) (service.Policy, error)
	GetOnGoingJobs() ([]JobInSchedule, error)
	StartJob(id uint) error
}

type Handler interface {
	Handle(JobInSchedule) error
}

type JobResult struct {
	JobID      uint
	Status     string
	ErrMessage string
}

func New(config *viper.Viper, db DB, handler Handler) *Scheduler {
	return &Scheduler{
		db:             db,
		newJobCh:       make(chan schedulerjob.Job),
		resultCh:       make(chan JobResult),
		DeleteBackupCh: make(chan backupset.Backupset),
		stoppingCh:     make(chan struct{}),
		handler:        handler,
		jobMutex:       make(map[string]struct{}),
	}
}

type Scheduler struct {
	db             DB
	newJobCh       chan schedulerjob.Job
	resultCh       chan JobResult
	DeleteBackupCh chan backupset.Backupset
	stoppingCh     chan struct{}
	handler        Handler
	config         *viper.Viper
	jobMutex       map[string]struct{}

	mu sync.Mutex
}

func (s *Scheduler) Start() {
	logger.Info("starting scheduler")

RunningLoop:
	for {
		select {
		case job := <-s.newJobCh:
			logger.Infof("received a new job: %v", job)
			if err := s.CheckMutex(job); err != nil {
				logger.Errorf("check mutex failed: %v", err)
				continue
			}

			if err := s.addJob(&job); err != nil {
				logger.Error(err)
			}
		case result := <-s.resultCh:
			logger.Infof("received result: %v", result)
			if err := s.completeJob(result.JobID, result.Status, result.ErrMessage); err != nil {
				logger.Error(err)
			}

			if err := s.Schedule(); err != nil {
				logger.Error(err)
			}
		case bset := <-s.DeleteBackupCh:
			s.ScheduleDelete(bset)
		case <-time.After(5 * time.Second):
			logger.Info("heart beat")
		case <-time.After(1 * time.Second):
			if err := s.Schedule(); err != nil {
				logger.Error(err)
			}
		case <-s.stoppingCh:
			break RunningLoop
		}
	}

	logger.Info("scheduler is stopped")
}

func (s *Scheduler) CheckMutex(job schedulerjob.Job) error {
	if job.JobType == schedulerjob.JobTypeBackupFile {
		key := fmt.Sprintf("%s-%d", job.JobType,job.PolicyID)

		if _, ok := s.jobMutex[key]; ok {
			return errors.New("job mutex occurs")
		}

		s.jobMutex[key] = struct{}{}
	}

	return nil
}

func (s *Scheduler) Stop() {
	logger.Info("stopping scheduler")
	s.stoppingCh <- struct{}{}
}

func (s *Scheduler) ScheduleDelete(bset backupset.Backupset) {
}

type BackupJobSpec struct {
}

func (s *Scheduler) addJob(job *schedulerjob.Job) error {
	now := time.Now()
	job.QueueTime = &now
	job.Status = schedulerjob.SchedulerJobStatusQueued
	return s.db.AddSchedulerJob(job)
}

func (s *Scheduler) completeJob(id uint, status string, msg string) error {
	return nil
}

func (s *Scheduler) Schedule() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobs, err := s.db.GetOnGoingJobs()

	if err != nil {
		return err
	}

	for _, job := range jobs {

		if job.Status == schedulerjob.SchedulerJobStatusQueued {
			s.StartJob(job)
		}
	}

	return nil
}

func (s *Scheduler) StartJob(job JobInSchedule) error {
	if err := s.db.StartJob(job.ID); err != nil {
		return err
	}

	s.handler.Handle(job)

	return nil
}

// DeleteBackup delete backup when repository is idle
func (s *Scheduler) DeleteBackup() error {
	return nil
}

func (s *Scheduler) AddSchedulerJobBackup(policyID uint, backupType string, operator string) error {
	var job schedulerjob.Job
	var jobType string

	policy, err := s.db.GetPolicy(policyID)

	if err != nil {
		return err
	}

	if policy.BackupSource.SourceType == "file" {
		jobType = schedulerjob.JobTypeBackupFile
	}

	job = schedulerjob.Job{
		PolicyID:     policy.ID,
		JobType:      jobType,
		BackupType:   backupType,
		Operator:     operator,
		Hostname:     policy.Hostname,
		Priority:     5,
		InSchedule:   true,
		RepositoryID: policy.RepositoryID,
		Status:       schedulerjob.SchedulerJobStatusCreated,
	}

	s.newJobCh <- job

	return nil
}
