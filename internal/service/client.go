package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/repository"
	"github.com/Coiiap5e/photographer/internal/utils"
)

type Client interface {
	CreateClient(ctx context.Context, client *model.Client) error
	DeleteClient(ctx context.Context, id int) error
	GetClients(ctx context.Context) error
	GetClientByID(ctx context.Context, id int) (*model.Client, error)
}

type postgresClient struct {
	clientRepo repository.Client
}

func NewClient(clientRepo repository.Client) Client {
	return &postgresClient{clientRepo: clientRepo}
}

func (c *postgresClient) CreateClient(ctx context.Context, client *model.Client) error {
	err := c.clientRepo.AddClient(ctx, client)
	if err != nil {
		return err
	}

	return nil
}

func (c *postgresClient) GetClientByID(ctx context.Context, id int) (*model.Client, error) {
	client, err := c.clientRepo.GetClientByID(ctx, id)
	if err != nil {
		if errors.IsErrorCode(err, errors.ErrCodeClientNotFound) {
			return nil, errors.Wrap(err, errors.ErrCodeClientNotFound, "client not found")
		}
		return nil, errors.Wrap(err, errors.ErrCodeDBSelect, "failed to get client")
	}

	return client, nil
}

func (c *postgresClient) DeleteClient(ctx context.Context, id int) error {
	for {
		confirm := utils.InputStringRequired("Are you sure you want to delete the client? (y/n)")
		if confirm == "n" || confirm == "N" {
			return errors.New(errors.ErrCodeValidation, "deletion cancelled")
		} else if confirm == "y" || confirm == "Y" {
			break
		} else {
			fmt.Println("Press wrong button: enter (y/n)")
		}
	}

	err := c.clientRepo.DeleteClient(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *postgresClient) GetClients(ctx context.Context) error {
	clients, err := c.clientRepo.GetClients(ctx)
	if err != nil {
		return err
	}

	showClients(clients)

	return nil
}

func showClients(clients []model.Client) {
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
