package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	cliapp "github.com/Coiiap5e/photographer/internal/app"
	"github.com/Coiiap5e/photographer/internal/config"
	"github.com/Coiiap5e/photographer/internal/database"
	"github.com/Coiiap5e/photographer/internal/repository"
	"github.com/Coiiap5e/photographer/internal/service"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	signalChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		logger.Error("configuration error", "error", err)
		os.Exit(1)
	}

	db, err := database.NewClient(ctx, dbConfig)
	if err != nil {
		logger.Error("error create db connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	clientRepo := repository.NewClient(db)
	shootRepo := repository.NewShoot(db)

	clientService := service.NewClient(clientRepo, logger)
	shootService := service.NewShoot(shootRepo, clientRepo, logger)

	go func() {
		sig := <-signalChan
		logger.Info("got signal", "signal", sig)

		done <- true
	}()

	app := cliapp.NewApp(clientService, shootService, logger)
	go func() {
		app.RunMenu(ctx)
		done <- true
	}()

	<-done
}
