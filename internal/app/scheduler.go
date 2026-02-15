package app

import (
	"context"
	"log/slog"

	"github.com/robfig/cron/v3"

	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
)

type Scheduler struct {
	logger       *slog.Logger
	shootService service.Shoot
	clock        *clock.Clock
	cron         *cron.Cron
}

func NewScheduler(logger *slog.Logger, shootService service.Shoot, clock *clock.Clock) *Scheduler {

	stdLogger := slog.NewLogLogger(logger.With("component", "cron").Handler(), slog.LevelInfo)

	c := cron.New(cron.WithLogger(cron.PrintfLogger(stdLogger)))
	return &Scheduler{
		logger:       logger,
		shootService: shootService,
		clock:        clock,
		cron:         c,
	}
}

func (s *Scheduler) Start() {
	s.logger.Info("starting cron scheduler")

	_, err := s.cron.AddFunc("0 * * * *", func() {
		s.logTodayShootsCount()
	})
	if err != nil {
		s.logger.Error("failed to schedule logTodayShootsCount", "error", err)
		return
	}

	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.logger.Info("stopping cron scheduler")
	s.cron.Stop()
}

func (s *Scheduler) logTodayShootsCount() {
	today := s.clock.Now()
	count, err := s.shootService.GetShootsCountByDate(context.Background(), today)
	if err != nil {
		s.logger.Error("failed to get today's shoots count", "error", err)
		return
	}
	s.logger.Info("today's shoots count", "count", count)
}
