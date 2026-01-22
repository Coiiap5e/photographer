package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
)

type Scheduler struct {
	logger       *slog.Logger
	shootService service.Shoot
	clock        *clock.Clock
}

func NewScheduler(logger *slog.Logger, shootService service.Shoot, clock *clock.Clock) *Scheduler {
	return &Scheduler{
		logger:       logger,
		shootService: shootService,
		clock:        clock,
	}
}

func (s *Scheduler) Start() {
	s.logger.Info("starting scheduler")
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.logTodayShootsCount()
			}
		}
	}()
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
