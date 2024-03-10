package scheduler

import (
	"github.com/spf13/viper"
	"rockbackup/backend/backupset"
	"sync"
	"time"
)

type DB interface {
}

type Handler interface {
	StartBackup(policyID uint) error
	StartRestore(backupsetID uint) error
}

type JobResult struct {
	JobID      uint
	Status     string
	ErrMessage string
}

func New(config *viper.Viper, db DB, handler Handler) *Scheduler {
	return &Scheduler{
		db:             db,
		newJobCh:       make(chan SchedulerJob),
		resultCh:       make(chan JobResult),
		DeleteBackupCh: make(chan backupset.Backupset),
		stoppingCh:     make(chan struct{}),
		handler:        handler,
	}
}

type Scheduler struct {
	db             DB
	newJobCh       chan SchedulerJob
	resultCh       chan JobResult
	DeleteBackupCh chan backupset.Backupset
	stoppingCh     chan struct{}
	handler        Handler
	config         *viper.Viper

	mu sync.Mutex
}

func (s *Scheduler) Start() {
	logger.Info("starting scheduler")

RunningLoop:
	for {
		select {
		case job := <-s.newJobCh:
			if err := s.AddJob(job); err != nil {
				logger.Error(err)
			}
		case result := <-s.resultCh:
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

func (s *Scheduler) Stop() {
	logger.Info("stopping scheduler")
	s.stoppingCh <- struct{}{}
}

func (s *Scheduler) ScheduleDelete(bset backupset.Backupset) {

}

func (s *Scheduler) AddJob(job SchedulerJob) error {
	return nil
}

func (s *Scheduler) completeJob(id uint, status string, msg string) error {
	return nil
}

func (s *Scheduler) Schedule() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

// DeleteBackup delete backup when repository is idle
func (s *Scheduler) DeleteBackup() error {
	return nil
}
