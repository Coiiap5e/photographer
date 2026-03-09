package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Coiiap5e/photographer/internal/adapter/currency"
	myerrors "github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/service"
)

type NotifyShootsJob struct {
	shootService service.Shoot
	notifier     service.Notifier
	logger       *slog.Logger
	window       time.Duration
	timeUnit     string
}

func NewNotifyShootsJob(shootService service.Shoot, notifier service.Notifier, logger *slog.Logger, window time.Duration) *NotifyShootsJob {
	unit := "hours"
	if window > 24*time.Hour {
		unit = "days"
	}
	return &NotifyShootsJob{
		shootService: shootService,
		notifier:     notifier,
		logger:       logger,
		window:       window,
		timeUnit:     unit,
	}
}

func (j *NotifyShootsJob) Execute() {
	var shoots []model.Shoot
	var err error

	ctx := context.Background()

	switch j.timeUnit {
	case "days":
		shoots, err = j.shootService.GetShootsForNextTwoDays(ctx)
	case "hours":
		shoots, err = j.shootService.GetShootsForNextTwoHours(ctx)
	default:
		err := myerrors.New(myerrors.ErrCodeJobError, fmt.Sprintf("unknown time unit for notification job: %s", j.timeUnit))
		j.logger.Error("job execution failed", "error", err)
		return
	}

	if err != nil {
		j.logger.Error(fmt.Sprintf("failed to get shoots for next %v", j.window), "error", err)
		return
	}

	if len(shoots) == 0 {
		message := fmt.Sprintf("No upcoming shoots in the next %.0f %s.", j.window.Hours(), j.timeUnit)
		j.logger.Info(message)
		if err := j.notifier.NotifyMessage(message); err != nil {
			j.logger.Error("failed to send notification", "error", err)
		}
		return
	}

	j.logger.Info(fmt.Sprintf("upcoming shoots found for next %.0f %s", j.window.Hours(), j.timeUnit), "count", len(shoots))
	for _, shoot := range shoots {
		if err := j.notifier.Notify(shoot); err != nil {
			j.logger.Error("failed to send notification for shoot", "shoot_id", shoot.Id, "error", err)
		}
	}
}

type UpdateCurrencyRateJob struct {
	realCurrencyService   currency.Service
	cachedCurrencyService *currency.InMemoryCacheService
	logger                *slog.Logger
}

func NewUpdateCurrencyRateJob(
	realCurrencyService currency.Service,
	cachedCurrencyService *currency.InMemoryCacheService,
	logger *slog.Logger,
) *UpdateCurrencyRateJob {
	return &UpdateCurrencyRateJob{
		realCurrencyService:   realCurrencyService,
		cachedCurrencyService: cachedCurrencyService,
		logger:                logger,
	}
}

func (j *UpdateCurrencyRateJob) Execute() {
	j.logger.Info("starting currency rate update job")

	rate, err := j.realCurrencyService.GetUSDRate(context.Background())
	if err != nil {
		j.logger.Error("failed to get USD rate from real service", "error", err)
		return
	}
	j.cachedCurrencyService.UpdateRate(rate)
	j.logger.Info("successfully updated cached currency rate", "new_rate", rate)
}
