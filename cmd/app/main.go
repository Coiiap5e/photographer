package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Coiiap5e/photographer/config"
	"github.com/Coiiap5e/photographer/internal/database"
	"github.com/Coiiap5e/photographer/internal/repository"
	"github.com/Coiiap5e/photographer/internal/service"
)

func main() {
	ctx := context.Background()

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
	shootService := service.NewShootService(shootRepo)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()
		fmt.Print("Select a menu item: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			err := clientService.CreateClient(ctx)
			if err != nil {
				fmt.Printf("Error creating client: %v\n", err)
			}
		case "2":
			err := clientService.DeleteClient(ctx)
			if err != nil {
				fmt.Printf("Error deleting client: %v\n", err)
			}
		case "3":
			err := shootService.CreateShoot(ctx)
			if err != nil {
				fmt.Printf("Error creating shoot: %v\n", err)
			}
		case "4":
			err := shootService.DeleteShoot(ctx)
			if err != nil {
				fmt.Printf("Error deleting shoot: %v\n", err)
			}
		case "5":
			err := clientService.GetClients(ctx)
			if err != nil {
				fmt.Printf("Error getting clients: %v\n", err)
			}
		case "6":
			err := shootService.GetShoots(ctx)
			if err != nil {
				fmt.Printf("Error getting shoots: %v\n", err)
			}
		case "7":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice")
		}
		fmt.Println("")
	}
}

func showMenu() {
	fmt.Println("1. Add client")
	fmt.Println("2. Delete client")
	fmt.Println("3. Add shoot")
	fmt.Println("4. Delete shoot")
	fmt.Println("5. Show list of clients")
	fmt.Println("6. Show list of shoots")
	fmt.Println("7. Exit")
}
