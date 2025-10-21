package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cliapp "github.com/Coiiap5e/photographer/internal/app"
	"github.com/Coiiap5e/photographer/internal/config"
	"github.com/Coiiap5e/photographer/internal/database"
	"github.com/Coiiap5e/photographer/internal/logs"
	"github.com/Coiiap5e/photographer/internal/repository"
	"github.com/Coiiap5e/photographer/internal/service"
)

func main() {
	ctx := context.Background()

	logger, closeLogger := logs.InitLogger()

	defer closeLogger()

	logger.Info("app starting")

	signalChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		fmt.Println("Db error! More information in logs")

		logger.Error("configuration error", "error", err)
		os.Exit(1)
	}

	db, err := database.NewClient(ctx, dbConfig)
	if err != nil {
		fmt.Println("Db error! More information in logs")

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

		fmt.Println("Got signal: ", sig)
		logger.Info("got signal", "signal", sig.String())

		done <- true
	}()

	app := cliapp.NewApp(clientService, shootService, logger)
	go func() {
		app.RunMenu(ctx)
		done <- true
	}()

	<-done

	logger.Info("app end")
}
