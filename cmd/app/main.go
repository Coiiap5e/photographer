package main

import (
	"context"
	"log"
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

	log.SetFlags(log.Ldate)
	log.SetOutput(os.Stdout)

	signalChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		log.Fatal("configuration error: ", err)
	}

	db, err := database.NewClient(ctx, dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	clientRepo := repository.NewClient(db)
	shootRepo := repository.NewShoot(db)

	clientService := service.NewClient(clientRepo)
	shootService := service.NewShoot(shootRepo, clientRepo)

	go func() {
		sig := <-signalChan
		log.Println("got signal:", sig)

		done <- true
	}()

	app := cliapp.NewApp(clientService, shootService)
	go func() {
		app.RunMenu(ctx)
		done <- true
	}()

	<-done
}
