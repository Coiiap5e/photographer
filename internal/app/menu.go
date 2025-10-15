package app

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/Coiiap5e/photographer/internal/service"
)

type App struct {
	clientService *service.ClientService
	shootService  *service.ShootService
}

func NewApp(clientService *service.ClientService, shootService *service.ShootService) *App {
	return &App{
		clientService: clientService,
		shootService:  shootService,
	}
}

func (a *App) RunMenu(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		a.showMenu()
		fmt.Print("Select a menu item: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			err := a.clientService.CreateClient(ctx)
			if err != nil {
				fmt.Printf("Error creating client: %v\n", err)
			}
		case "2":
			err := a.clientService.DeleteClient(ctx)
			if err != nil {
				fmt.Printf("Error deleting client: %v\n", err)
			}
		case "3":
			err := a.shootService.CreateShoot(ctx)
			if err != nil {
				fmt.Printf("Error creating shoot: %v\n", err)
			}
		case "4":
			err := a.shootService.DeleteShoot(ctx)
			if err != nil {
				fmt.Printf("Error deleting shoot: %v\n", err)
			}
		case "5":
			err := a.clientService.GetClients(ctx)
			if err != nil {
				fmt.Printf("Error getting clients: %v\n", err)
			}
		case "6":
			err := a.shootService.GetShoots(ctx)
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

func (a *App) showMenu() {
	fmt.Println("1. Add client")
	fmt.Println("2. Delete client")
	fmt.Println("3. Add shoot")
	fmt.Println("4. Delete shoot")
	fmt.Println("5. Show list of clients")
	fmt.Println("6. Show list of shoots")
	fmt.Println("7. Exit")
}
