package repository

import (
	"context"
	"errors"

	"github.com/Coiiap5e/photographer/internal/database"
	myerrors "github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/model"
	"github.com/jackc/pgx/v5"
)

type ClientRepository struct {
	db *database.DB
}

func NewClientRepository(db *database.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (repo *ClientRepository) AddClient(ctx context.Context, client *model.Client) error {
	query := `
INSERT INTO clients 
    (first_name, last_name, phone, social_network_url) 
VALUES 
    ($1, $2, $3, $4) 
RETURNING 
    id, created_at`

	err := repo.db.Pool.QueryRow(ctx, query,
		client.FirstName, client.LastName, client.Phone, client.SocialNetworkUrl).
		Scan(&client.Id, &client.CreatedAt)

	if err != nil {
		return myerrors.Wrap(err, myerrors.ErrCodeClientCreate, "failed to create client")
	}

	return nil
}

func (repo *ClientRepository) DeleteClient(ctx context.Context, id int) error {
	query := `
DELETE FROM clients WHERE id = $1`

	result, err := repo.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return myerrors.Wrap(err, myerrors.ErrCodeClientDelete, "failed to delete client")
	}

	if result.RowsAffected() == 0 {
		return myerrors.New(myerrors.ErrCodeClientNotFound, "client not found")
	}

	return nil
}

func (repo *ClientRepository) GetClientByID(ctx context.Context, id int) (*model.Client, error) {
	query := `
SELECT id, first_name, last_name, phone, social_network_url
FROM clients
WHERE id = $1`

	var client model.Client
	err := repo.db.Pool.QueryRow(ctx, query, id).Scan(
		&client.Id, &client.FirstName, &client.LastName,
		&client.Phone, &client.SocialNetworkUrl)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerrors.New(myerrors.ErrCodeClientNotFound, "client not found")
		}
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get client")
	}

	return &client, nil
}

func (repo *ClientRepository) GetClients(ctx context.Context) ([]model.Client, error) {
	query := `
SELECT id, first_name, last_name, phone, social_network_url, created_at
FROM clients`

	rows, err := repo.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get clients")
	}

	defer rows.Close()

	var clients []model.Client

	for rows.Next() {
		var client model.Client
		err := rows.Scan(
			&client.Id, &client.FirstName, &client.LastName,
			&client.Phone, &client.SocialNetworkUrl, &client.CreatedAt)
		if err != nil {
			return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get client")
		}
		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "error during rows iteration")
	}

	return clients, nil
}
