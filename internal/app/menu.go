package app

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/utils"
)

type App struct {
	clientService service.Client
	shootService  service.Shoot
}

func NewApp(clientService service.Client, shootService service.Shoot) *App {
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
			client := &model.Client{
				FirstName:        utils.InputStringRequired("First name"),
				LastName:         utils.InputStringRequired("Last name"),
				Phone:            utils.InputStringRequired("Phone number"),
				SocialNetworkUrl: utils.InputString("Social network url"),
			}
			err := a.clientService.CreateClient(ctx, client)

			if err != nil {
				fmt.Printf("Error creating client: %v\n", err)
			}

			log.Printf("client added with ID: %d, created at: %v",
				client.Id, client.CreatedAt)

		case "2":
			var id int
			var client *model.Client
			var err error

			for {
				id = utils.InputId("ID of the client")
				client, err = a.clientService.GetClientByID(ctx, id)

				if err != nil {
					if errors.IsErrorCode(err, errors.ErrCodeClientNotFound) {
						fmt.Println("Client not found. Try again")
						continue
					} else {
						fmt.Printf("Failed to get client: %v\n", err)
					}
				}

				break
			}

			fmt.Printf("Confirm deleting client: %s %s\n", client.FirstName, client.LastName)
			fmt.Printf("Phone number: %s\n", client.Phone)
			if client.SocialNetworkUrl != "" {
				fmt.Printf("Social network url: %s\n", client.SocialNetworkUrl)
			}

			err = a.clientService.DeleteClient(ctx, id)
			if err != nil {
				fmt.Printf("Error deleting client: %v\n", err)
			}

			log.Printf("client %s %s with ID: %d deleted successfully \n",
				client.FirstName, client.LastName, id)

		case "3":
			var shoot *model.Shoot

			var clientId int
			var clientFirstName, clientLastName string

			shootDate, startTime, endTime := utils.InputShootDate()

			for {
				client, err := a.clientService.GetClientByID(ctx, utils.InputId("Client_id"))
				if err != nil {
					if errors.IsErrorCode(err, errors.ErrCodeClientNotFound) {
						fmt.Println("Client not found. Try again")
						continue
					} else {
						fmt.Printf("Failed to get client: %v\n", err)
					}
				}
				clientId = client.Id
				clientFirstName = client.FirstName
				clientLastName = client.LastName
				break
			}

			shoot = &model.Shoot{
				ClientId:        clientId,
				ShootDate:       shootDate,
				StartTime:       startTime,
				EndTime:         endTime,
				ShootPrice:      utils.InputInt("Shoot price"),
				ShootLocation:   utils.InputStringRequired("Location"),
				ClientFirstName: clientFirstName,
				ClientLastName:  clientLastName,
				ShootType:       utils.InputStringRequired("Shoot type"),
				Notes:           utils.InputString("Notes"),
			}

			err := a.shootService.CreateShoot(ctx, shoot)
			if err != nil {
				fmt.Printf("Error creating shoot: %v\n", err)
			}

			log.Printf("shoot added successfully")

		case "4":
			var id int
			var shoot *model.Shoot
			var err error

			for {
				id = utils.InputId("ID of the shoot")
				shoot, err = a.shootService.GetShootByID(ctx, id)

				if err != nil {
					if errors.IsErrorCode(err, errors.ErrCodeShootNotFound) {
						fmt.Println("Shoot not found. Try again")
						continue
					} else {
						fmt.Printf("Failed to get shoot: %v\n", err)
					}
				}

				break
			}

			fmt.Printf("Confirm deleting shoot: %s start: %s end: %s\n",
				shoot.ShootDate.Format("02.01.2006"),
				shoot.StartTime.Format("15:04"),
				shoot.EndTime.Format("15:04"))
			fmt.Printf("Location: %s. ShootType: %s. Price: %d\n",
				shoot.ShootLocation, shoot.ShootType, shoot.ShootPrice)
			fmt.Printf("Client id: %d name: %s %s\n", shoot.ClientId,
				shoot.ClientFirstName, shoot.ClientLastName)
			if shoot.Notes != "" {
				fmt.Printf("Notes: %s\n", shoot.Notes)
			}

			err = a.shootService.DeleteShoot(ctx, id)
			if err != nil {
				fmt.Printf("Error deleting shoot: %v\n", err)
			}

			log.Printf("shoot %s with %s %s deleted successfully \n",
				shoot.StartTime.Format("02.01.2006 15:04"),
				shoot.ClientFirstName, shoot.ClientLastName)

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
