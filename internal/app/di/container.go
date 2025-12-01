package di

import (
	"context"
	"log/slog"

	"github.com/Coiiap5e/photographer/internal/api/controllers"
	"github.com/Coiiap5e/photographer/internal/config"
	"github.com/Coiiap5e/photographer/internal/database"
	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/handlers"
	"github.com/Coiiap5e/photographer/internal/logs"
	"github.com/Coiiap5e/photographer/internal/repository"
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
	Handlers     *Handlers

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

type Handlers struct {
	Shoot  *handlers.ShootHandler
	Client *handlers.ClientHandler
}

func NewContainer(ctx context.Context) (*Container, error) {
	newClock := clock.NewInMoscow()

	logger, closeLogger := logs.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		closeLogger()
		return nil, errors.Wrap(err, errors.ErrCodeDBConfig, "configuration error")
	}

	db, err := database.NewClient(ctx, cfg.DB)
	if err != nil {
		closeLogger()
		return nil, errors.Wrap(err, errors.ErrCodeDBConnection, "error create db connection")
	}

	clientRepo := repository.NewClient(db, newClock)
	shootRepo := repository.NewShoot(db, newClock)

	clientService := service.NewClient(clientRepo, logger)
	shootService := service.NewShoot(shootRepo, clientService, logger)

	clientController := controllers.NewClientController(clientService)
	shootController := controllers.NewShootController(shootService, clientService, newClock)

	clientHandler := handlers.NewClientHandler(clientController)
	shootHandler := handlers.NewShootHandler(shootController)

	return &Container{
		Config:      cfg,
		DB:          db,
		Clock:       newClock,
		Logger:      logger,
		closeLogger: closeLogger,
		Services: &Services{
			Shoot:  shootService,
			Client: clientService,
		},
		Repositories: &Repositories{
			Shoot:  shootRepo,
			Client: clientRepo,
		},
		Controllers: &Controllers{
			Shoot:  shootController,
			Client: clientController,
		},
		Handlers: &Handlers{
			Shoot:  shootHandler,
			Client: clientHandler,
		},
	}, nil

}

func (c *Container) Close() {
	c.DB.Close()

	if c.closeLogger != nil {
		c.closeLogger()
	}
}
