package worker

import (
	"log/slog"
	"sync"
)

// Job represents a job to be executed.
type Job interface {
	Execute()
}

// Pool is a worker pool that executes jobs concurrently.
type Pool struct {
	jobQueue chan Job
	wg       sync.WaitGroup
	workers  int
	logger   *slog.Logger
}

// NewPool creates a new worker pool.
func NewPool(workers int, queueSize int, logger *slog.Logger) *Pool {
	return &Pool{
		jobQueue: make(chan Job, queueSize),
		workers:  workers,
		logger:   logger,
	}
}

// Start starts the worker pool.
func (p *Pool) Start() {
	p.logger.Info("starting worker pool", "workers", p.workers)
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			p.logger.Info("worker started", "worker_id", workerID)
			for job := range p.jobQueue {
				p.logger.Info("worker received job", "worker_id", workerID)
				job.Execute()
				p.logger.Info("worker finished job", "worker_id", workerID)
			}
			p.logger.Info("worker stopped", "worker_id", workerID)
		}(i + 1)
	}
}

// Submit submits a job to the worker pool.
func (p *Pool) Submit(job Job) {
	p.logger.Info("submitting job")
	p.jobQueue <- job
}

// Stop stops the worker pool gracefully.
func (p *Pool) Stop() {
	p.logger.Info("stopping worker pool")
	close(p.jobQueue)
	p.wg.Wait()
	p.logger.Info("worker pool stopped")
}
