package di

import (
	"context"
	"log/slog"

	"github.com/Coiiap5e/photographer/internal/adapter/repository"
	"github.com/Coiiap5e/photographer/internal/api/controllers"
	"github.com/Coiiap5e/photographer/internal/app"
	"github.com/Coiiap5e/photographer/internal/config"
	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/infrastructure/database"
	"github.com/Coiiap5e/photographer/internal/logs"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
)

type Container struct {
	Config       *config.Config
	DB           *database.DB
	Clock        *clock.Clock
	Logger       *slog.Logger
	Services     *Services
	Repositories *Repositories
	Controllers  *Controllers
	Scheduler    *app.Scheduler

	closeLogger func()
}

type Services struct {
	Shoot  service.Shoot
	Client service.Client
}

type Repositories struct {
	Shoot  repository.Shoot
	Client repository.Client
}

type Controllers struct {
	Shoot  *controllers.ShootController
	Client *controllers.ClientController
}

func NewContainer(ctx context.Context) (*Container, error) {
	container := &Container{}

	container.Clock = clock.NewInMoscow()

	container.Logger, container.closeLogger = logs.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		container.closeLogger()
		return nil, errors.Wrap(err, errors.ErrCodeDBConfig, "configuration error")
	}
	container.Config = cfg

	db, err := database.NewClient(ctx, cfg.DB)
	if err != nil {
		container.closeLogger()
		return nil, errors.Wrap(err, errors.ErrCodeDBConnection, "error create db connection")
	}
	container.DB = db

	clientRepo := repository.NewClient(db, container.Clock)
	shootRepo := repository.NewShoot(db, container.Clock)

	container.Repositories = &Repositories{
		Client: clientRepo,
		Shoot:  shootRepo,
	}

	clientService := service.NewClient(clientRepo, container.Logger)
	shootService := service.NewShoot(shootRepo, clientService, container.Logger)

	container.Services = &Services{
		Client: clientService,
		Shoot:  shootService,
	}

	clientController := controllers.NewClientController(container.Services.Client)
	shootController := controllers.NewShootController(container.Services.Shoot, container.Services.Client, container.Clock)
	container.Controllers = &Controllers{
		Client: clientController,
		Shoot:  shootController,
	}

	container.Scheduler = app.NewScheduler(container.Logger, container.Services.Shoot, container.Clock)

	return container, nil
}

func (c *Container) Close() {
	c.DB.Close()

	if c.closeLogger != nil {
		c.closeLogger()
	}
}
