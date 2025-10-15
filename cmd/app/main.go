package main

import (
	"context"
	"log"
	"os"

	"github.com/Coiiap5e/photographer/config"
	cliapp "github.com/Coiiap5e/photographer/internal/app"
	"github.com/Coiiap5e/photographer/internal/database"
	"github.com/Coiiap5e/photographer/internal/repository"
	"github.com/Coiiap5e/photographer/internal/service"
)

func main() {
	ctx := context.Background()

	log.SetFlags(log.Ldate)
	log.SetOutput(os.Stdout)

	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		log.Fatal("Configuration error: ", err)
	}

	db, err := database.NewClient(ctx, dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	clientRepo := repository.NewClientRepository(db)
	clientService := service.NewClientService(clientRepo)
	shootRepo := repository.NewShootRepository(db)
	shootService := service.NewShootService(shootRepo, clientRepo)

	app := cliapp.NewApp(clientService, shootService)
	app.RunMenu(ctx)
}
