package worker

import (
	"context"
	"rockbackup/backend/agentd"
	"rockbackup/backend/policy"
	"rockbackup/backend/repository"
	"rockbackup/backend/schedulerjob"
	"sync"
)

// Task represents a backup task
type Task struct {
	Agent      *agentd.Agent
	BackupType string
	Instance   policy.Instance
	Repo       repository.Repository
	TargetPath string
}

type JobPool struct {
	jobChan  chan Task
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	maxprocs uint
}

// NewWorker creates a new Worker instance
func NewWorker() *JobPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &JobPool{
		taskChan: make(chan Task),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start starts the worker to process tasks
func (p *JobPool) Start() {
	for i := uint(0); i < p.maxprocs; i++ {
		go func() error {
			p.worker(p.jobChan)
			return nil
		}()
	}
}

// AddTask adds a new task to the worker's queue
func (p *JobPool) AddTask(task Task) {
	p.taskChan <- task
}

// processTask processes a single backup task
func (w *JobPool) worker(ch chan schedulerjob.Job) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok = <-jobs:
			if !ok {
				return
			}
		}

	}
}
