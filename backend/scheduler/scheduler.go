package scheduler

import (
	"errors"
	"fmt"
	"rockbackup/backend/backupset"
	"rockbackup/backend/log"
	"rockbackup/backend/policy"
	"rockbackup/backend/repository"
	"rockbackup/backend/schedulerjob"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var logger *log.Logger

func init() {
	logName := "job-scheduler"
	logger = log.New(logName)
}

type JobInSchedule struct {
	schedulerjob.Job
	Repository repository.Repository
}

type DB interface {
	AddSchedulerJob(*schedulerjob.Job) error
	GetPolicy(uint) (policy.Policy, error)
	GetJobsInschedule() ([]JobInSchedule, error)
	StartJob(id uint) error
	GetBackupset(uint) (backupset.Backupset, error)
}

type JobHandler interface {
	Start(JobInSchedule) error
}

type JobResult struct {
	JobID      uint
	Status     string
	ErrMessage string
}

func New(config *viper.Viper, db DB, handler JobHandler) *Scheduler {

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
	handler        JobHandler
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

			if err := s.addJob(&job); err != nil {
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
		key := fmt.Sprintf("%s-%d", job.JobType, job.PolicyID)

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

// func (s *Scheduler) completeJob(id uint, status string, msg string) error {
// 	s.db.CompleteJob()
// 	key := fmt.Sprintf("%s-%d", job.JobType, job.PolicyID)
// 	delete(s.jobMutex, key)
// 	return nil
// }

func (s *Scheduler) Schedule() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// logger.Info("Scheduling jobs")

	jobs, err := s.db.GetJobsInschedule()

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

	s.handler.Start(job)

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

func (s *Scheduler) AddSchedulerJobRestore(policyID uint, backupsetID uint, targetPath string, operator string) error {
	var job schedulerjob.Job
	var jobType string

	policy, err := s.db.GetPolicy(policyID)

	backupset, err := s.db.GetBackupset(backupsetID)

	if err != nil {
		return err
	}

	if policy.BackupSource.SourceType == "file" {
		jobType = schedulerjob.JobTypeBackupFile
	}

	job = schedulerjob.Job{
		PolicyID:     policy.ID,
		JobType:      jobType,
		Operator:     operator,
		Hostname:     policy.Hostname,
		Priority:     5,
		InSchedule:   true,
		RepositoryID: backupset.RepositoryID,
		Status:       schedulerjob.SchedulerJobStatusCreated,
	}

	s.newJobCh <- job

	return nil
}
