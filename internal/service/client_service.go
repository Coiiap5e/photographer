package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/repository"
	"github.com/Coiiap5e/photographer/internal/utils"
	"github.com/Coiiap5e/photographer/model"
)

type ClientService struct {
	clientRepo *repository.ClientRepository
}

func NewClientService(clientRepo *repository.ClientRepository) *ClientService {
	return &ClientService{clientRepo: clientRepo}
}

func (c *ClientService) CreateClient(ctx context.Context) error {
	client := &model.Client{
		FirstName:        utils.InputStringRequired("First name"),
		LastName:         utils.InputStringRequired("Last name"),
		Phone:            utils.InputStringRequired("Phone number"),
		SocialNetworkUrl: utils.InputString("Social network url"),
	}

	err := c.clientRepo.AddClient(ctx, client)
	if err != nil {
		return err
	}

	log.Printf("Client added with ID: %d, created at: %v", client.Id, client.CreatedAt)
	return nil
}

func (c *ClientService) DeleteClient(ctx context.Context) error {
	id := utils.InputId("ID of the client")

	client, err := c.clientRepo.GetClientByID(ctx, id)
	if err != nil {
		return err
	}

	fmt.Println("Confirm the deletion")
	fmt.Printf("Deleting client: %s %s\n", client.FirstName, client.LastName)
	fmt.Printf("Phone number: %s\n", client.Phone)
	if client.SocialNetworkUrl != "" {
		fmt.Printf("Social network url: %s\n", client.SocialNetworkUrl)
	}

	confirm := utils.InputStringRequired("Are you sure you want to delete the client? (y/n)")
	if confirm == "n" || confirm == "N" {
		return errors.New(errors.ErrCodeValidation, "deletion cancelled")
	}

	err = c.clientRepo.DeleteClient(ctx, id)
	if err != nil {
		return err
	}

	log.Printf("Client %s %s with ID: %d deleted successfully \n",
		client.FirstName, client.LastName, id)
	return nil
}

func (c *ClientService) GetClients(ctx context.Context) error {
	clients, err := c.clientRepo.GetClients(ctx)
	if err != nil {
		return err
	}

	c.showClients(clients)

	return nil
}

func (c *ClientService) showClients(clients []model.Client) {
	if len(clients) == 0 {
		fmt.Println("No clients found")
		return
	}

	fmt.Printf("%-4s %-15s %-15s %-20s %-25s %-12s\n",
		"ID", "First Name", "Last Name", "Phone", "Social Network", "Created")
	fmt.Println(strings.Repeat("-", 95))

	for _, client := range clients {
		fmt.Printf("%-4d %-15s %-15s %-20s %-25s %-12s\n",
			client.Id,
			client.FirstName,
			client.LastName,
			client.Phone,
			client.SocialNetworkUrl,
			client.CreatedAt.Format("02.01.2006"),
		)
	}
}
