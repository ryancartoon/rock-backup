package pool

import (
	"context"
	"log"
	"rockbackup/backend/schedulerjob"
	"sync"
)

type JobStarter interface {
	StartJob(schedulerjob.Job) error
}


type JobPool struct {
	jobChan  chan schedulerjob.Job
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	maxprocs uint
	JobStarter JobStarter
}

// NewWorker creates a new Worker instance
func NewWorker() *JobPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &JobPool{
		jobChan: make(chan schedulerjob.Job),
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
func (p *JobPool) AddTask(job schedulerjob.Job) {
	p.jobChan <- job
}

// processTask processes a single backup task
func (w *JobPool) worker(ch chan schedulerjob.Job) {
	for {
		select {
		case job, ok := <-ch:
			if !ok {
				return
			}

			err := w.JobStarter.StartJob(job)
			if err != nil {
				log.Printf("Error starting job: %v", err)
			}
		}

	}
}
