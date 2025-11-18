//go:build wireinject
// +build wireinject

package main

import (
	"time"

	cliapp "github.com/Coiiap5e/photographer/internal/app"
	"github.com/Coiiap5e/photographer/internal/config"
	"github.com/Coiiap5e/photographer/internal/database"
	"github.com/Coiiap5e/photographer/internal/logs"
	"github.com/Coiiap5e/photographer/internal/repository"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/google/wire"
)

//go:generate wire

func InitializeApp() (*cliapp.App, func(), error) {
	wire.Build(
		provideDBConfig,
		providePoolConfig,

		logs.InitLogger,

		database.NewClient,

		repository.NewClient,
		repository.NewShoot,

		service.NewClient,
		service.NewShoot,

		cliapp.NewApp,
	)

	return nil, nil, nil
}

func provideDBConfig() (config.DbConfig, error) {
	return config.LoadDBConfig()
}

func providePoolConfig() config.PoolConfig {
	return config.NewPoolConfigBuilder().
		WithMaxOpenConns(25).
		WithMaxIdleConns(5).
		WithMaxConnLifetime(30 * time.Minute).
		WithMaxConnIdleTime(5 * time.Minute).
		Build()
}
