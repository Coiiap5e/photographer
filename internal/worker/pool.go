package worker

import (
	"context"
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

// Run starts the worker pool and blocks until the context is cancelled.
func (p *Pool) Run(ctx context.Context) error {
	p.logger.Info("starting worker pool", "workers", p.workers)

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			p.logger.Info("worker started", "worker_id", workerID)
			for {
				select {
				case job, ok := <-p.jobQueue:
					if !ok {
						p.logger.Info("worker stopped (job queue closed)", "worker_id", workerID)
						return
					}
					p.logger.Info("worker received job", "worker_id", workerID)
					job.Execute()
					p.logger.Info("worker finished job", "worker_id", workerID)
				case <-ctx.Done():
					p.logger.Info("worker stopped (context cancelled)", "worker_id", workerID)
					return
				}
			}
		}(i + 1)
	}

	// Wait for context cancellation to initiate shutdown.
	<-ctx.Done()

	p.logger.Info("stopping worker pool (context cancelled)")
	close(p.jobQueue)
	p.wg.Wait()
	p.logger.Info("worker pool stopped")

	return nil
}

// Submit submits a job to the worker pool.
func (p *Pool) Submit(job Job) {
	p.logger.Info("submitting job")
	p.jobQueue <- job
}
