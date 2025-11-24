package app

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/utils"
	"github.com/samber/lo"
)

type App struct {
	clientService service.Client
	shootService  service.Shoot
	logger        *slog.Logger
}

func NewApp(clientService service.Client, shootService service.Shoot, logger *slog.Logger) *App {
	return &App{
		clientService: clientService,
		shootService:  shootService,
		logger:        logger,
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

			fmt.Println("Client added")
			a.logger.Info("client added",
				"id", client.Id,
				"created at", client.CreatedAt,
			)

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

			fmt.Println("Client deleted successfully")
			a.logger.Info("client deleted successfully",
				"client first name", client.FirstName,
				"client last name", client.LastName,
				"client id", id,
			)

		case "3":
			clients := make([]model.ShootClient, 0)

			shootDate, startTime, endTime := utils.InputShootDate()

			for {
				clientID := utils.InputId("Client ID")

				if lo.ContainsBy(clients, func(checkClient model.ShootClient) bool {
					return checkClient.ClientID == clientID
				}) {
					fmt.Printf("Client with ID %d is already added to this shoot\n", clientID)
					continue
				}

				client, err := a.clientService.GetClientByID(ctx, clientID)
				if err != nil {
					if errors.IsErrorCode(err, errors.ErrCodeClientNotFound) {
						fmt.Println("Client not found. Try again")
						continue
					} else {
						fmt.Printf("Failed to get client: %v\n", err)
					}
				}

				isMainClient := utils.InputBool("Is this main client? (true/false)")
				relationshipType := utils.InputStringRequired("Relationship type (bride/groom/witness/etc)")

				shootClient := model.ShootClient{
					ClientID:         client.Id,
					IsMainClient:     isMainClient,
					RelationshipType: relationshipType,
				}

				clients = append(clients, shootClient)

				var confirm string
				for {
					confirm = utils.InputStringRequired("Add another client? (y/n)")
					confirm = strings.ToLower(confirm)
					if confirm == "y" || confirm == "n" {
						break
					}
					fmt.Println("Press wrong button: enter (y/n)")
				}

				if confirm == "n" {
					break
				}
			}

			shoot := &model.Shoot{
				ShootDate:     shootDate,
				StartTime:     startTime,
				EndTime:       endTime,
				ShootPrice:    utils.InputInt("Shoot price"),
				ShootLocation: utils.InputStringRequired("Location"),
				ShootType:     utils.InputStringRequired("Shoot type"),
				Notes:         utils.InputString("Notes"),
			}

			err := a.shootService.CreateShoot(ctx, shoot, clients)
			if err != nil {
				fmt.Printf("Error creating shoot: %v\n", err)
			}

			fmt.Println("shoot added successfully")
			a.logger.Info("shoot added successfully")

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

			clients := shoot.Clients

			fmt.Printf("Confirm deleting shoot: %s start: %s end: %s\n",
				shoot.ShootDate.Format("02.01.2006"),
				shoot.StartTime.Format("15:04"),
				shoot.EndTime.Format("15:04"))
			fmt.Printf("Location: %s. ShootType: %s. Price: %d\n",
				shoot.ShootLocation, shoot.ShootType, shoot.ShootPrice)
			fmt.Printf("Clients:")
			for index, client := range clients {
				fmt.Printf("Client №%d: id: %d name: %s %s\n", index+1, client.ClientID,
					client.FirstName, client.LastName)
			}
			if shoot.Notes != "" {
				fmt.Printf("Notes: %s\n", shoot.Notes)
			}

			err = a.shootService.DeleteShoot(ctx, id)
			if err != nil {
				fmt.Printf("Error deleting shoot: %v\n", err)
			}

			fmt.Println("shoot deleted successfully")
			a.logger.Info("shoot deleted successfully",
				"start date", shoot.StartTime.Format("02.01.2006 15:04"),
			)

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
			err := a.shootService.GetShootsSortedByDate(ctx)
			if err != nil {
				fmt.Printf("Error getting shoots: %v\n", err)
			}
		case "8":
			err := a.shootService.GetShootsWithRelationshipType(ctx, "child")
			if err != nil {
				fmt.Printf("Error getting shoots: %v\n", err)
			}
		case "9":
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
	fmt.Println("7. Calendar of shoots")
	fmt.Println("8. Show shoots with childs")
	fmt.Println("9. Exit")
}
