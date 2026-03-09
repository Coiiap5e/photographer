package app

import (
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/Coiiap5e/photographer/internal/adapter/currency"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
	"github.com/Coiiap5e/photographer/internal/worker"
)

type Scheduler struct {
	logger                *slog.Logger
	shootService          service.Shoot
	realCurrencyService   currency.Service
	cachedCurrencyService *currency.InMemoryCacheService
	notifier              service.Notifier
	clock                 *clock.Clock
	cron                  *cron.Cron
	pool                  *worker.Pool
}

func NewScheduler(
	logger *slog.Logger,
	shootService service.Shoot,
	realCurrencyService currency.Service,
	cachedCurrencyService *currency.InMemoryCacheService,
	notifier service.Notifier,
	clock *clock.Clock,
	pool *worker.Pool,
) *Scheduler {

	stdLogger := slog.NewLogLogger(logger.With("component", "cron").Handler(), slog.LevelInfo)
	c := cron.New(cron.WithLogger(cron.PrintfLogger(stdLogger)))
	return &Scheduler{
		logger:                logger,
		shootService:          shootService,
		realCurrencyService:   realCurrencyService,
		cachedCurrencyService: cachedCurrencyService,
		notifier:              notifier,
		clock:                 clock,
		cron:                  c,
		pool:                  pool,
	}
}

func (s *Scheduler) Start() {
	s.logger.Info("starting cron scheduler")

	// Schedule shoot notifications
	_, err := s.cron.AddFunc("0 9 * * *", func() {
		job := NewNotifyShootsJob(s.shootService, s.notifier, s.logger, 48*time.Hour)
		s.pool.Submit(job)
	})
	if err != nil {
		s.logger.Error("failed to schedule two-day notification job", "error", err)
		return
	}

	_, err = s.cron.AddFunc("0 * * * *", func() {
		job := NewNotifyShootsJob(s.shootService, s.notifier, s.logger, 2*time.Hour)
		s.pool.Submit(job)
	})
	if err != nil {
		s.logger.Error("failed to schedule two-hour notification job", "error", err)
		return
	}

	// Schedule currency rate updates
	_, err = s.cron.AddFunc("0 * * * *", func() { // Every hour
		job := NewUpdateCurrencyRateJob(s.realCurrencyService, s.cachedCurrencyService, s.logger)
		s.pool.Submit(job)
	})
	if err != nil {
		s.logger.Error("failed to schedule currency rate update job", "error", err)
		return
	}

	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.logger.Info("stopping cron scheduler")
	s.cron.Stop()
}

