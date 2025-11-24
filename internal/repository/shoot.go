package repository

import (
	"context"
	"errors"

	"github.com/Coiiap5e/photographer/internal/database"
	myerrors "github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
	"github.com/jackc/pgx/v5"
)

type Shoot interface {
	AddShoot(ctx context.Context, shoot *model.Shoot, clients []model.ShootClient) error
	DeleteShoot(ctx context.Context, id int) error
	GetShootByID(ctx context.Context, id int) (*model.Shoot, error)
	GetShoots(ctx context.Context) ([]model.Shoot, error)
}

type postgresShoot struct {
	db    *database.DB
	clock *clock.Clock
}

func NewShoot(db *database.DB, clock *clock.Clock) Shoot {
	return &postgresShoot{
		db:    db,
		clock: clock,
	}
}

func (repo *postgresShoot) AddShoot(ctx context.Context, shoot *model.Shoot, clients []model.ShootClient) error {
	tx, err := repo.db.Pool.Begin(ctx)
	if err != nil {
		return myerrors.Wrap(err, myerrors.ErrCodeDBTransaction, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	query := `
INSERT INTO shoots
	(date, start_time, end_time, shoot_price, location, shoot_type, notes, created_at)
VALUES 
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
	id`

	err = tx.QueryRow(ctx, query, shoot.ShootDate,
		shoot.StartTime, shoot.EndTime, shoot.ShootPrice, shoot.ShootLocation,
		shoot.ShootType, shoot.Notes, repo.clock.Now()).Scan(&shoot.Id)

	if err != nil {
		return myerrors.Wrap(err, myerrors.ErrCodeShootCreate, "failed to create shoot")
	}

	for i := range clients {
		clients[i].ShootID = shoot.Id

		err := repo.addShootClientTx(ctx, tx, clients[i])
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return myerrors.Wrap(err, myerrors.ErrCodeDBTransaction, "failed to commit transaction")
	}

	return nil
}

func (repo *postgresShoot) addShootClientTx(ctx context.Context, tx pgx.Tx, client model.ShootClient) error {
	query := `
INSERT INTO shoot_clients 
    (shoot_id, client_id, is_main_client, relationship_type, created_at)
VALUES
	($1, $2, $3, $4, $5)
`
	_, err := tx.Exec(ctx, query, client.ShootID, client.ClientID, client.IsMainClient, client.RelationshipType, repo.clock.Now())
	if err != nil {
		return myerrors.Wrap(err, myerrors.ErrCodeDBQuery, "failed to add client to shoot")
	}
	return nil
}

func (repo *postgresShoot) DeleteShoot(ctx context.Context, id int) error {
	query := `
DELETE FROM shoots WHERE id = $1`

	result, err := repo.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return myerrors.Wrap(err, myerrors.ErrCodeShootDelete, "failed to delete shoot")
	}

	if result.RowsAffected() == 0 {
		return myerrors.New(myerrors.ErrCodeShootNotFound, "shoot not found")
	}

	return nil
}

func (repo *postgresShoot) GetShootByID(ctx context.Context, id int) (*model.Shoot, error) {
	query := `
SELECT 
	id, date, start_time, end_time, 
	shoot_price, location, shoot_type, notes
FROM shoots
WHERE id = $1`

	var shoot model.Shoot
	err := repo.db.Pool.QueryRow(ctx, query, id).Scan(
		&shoot.Id, &shoot.ShootDate, &shoot.StartTime,
		&shoot.EndTime, &shoot.ShootPrice, &shoot.ShootLocation,
		&shoot.ShootType, &shoot.Notes)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerrors.New(myerrors.ErrCodeShootNotFound, "shoot not found")
		}
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoot")
	}

	clients, err := repo.getShootClients(ctx, id)
	if err != nil {
		return nil, err
	}

	shoot.Clients = clients
	return &shoot, nil
}

func (repo *postgresShoot) GetShoots(ctx context.Context) ([]model.Shoot, error) {
	query := `
SELECT 
    id, date, start_time, end_time, 
    shoot_price, location, shoot_type, notes, created_at
FROM shoots
ORDER BY date DESC, created_at DESC
`

	rows, err := repo.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoots")
	}

	defer rows.Close()

	shoots := make([]model.Shoot, 0)

	for rows.Next() {
		var shoot model.Shoot
		err := rows.Scan(
			&shoot.Id, &shoot.ShootDate, &shoot.StartTime,
			&shoot.EndTime, &shoot.ShootPrice, &shoot.ShootLocation,
			&shoot.ShootType, &shoot.Notes, &shoot.CreatedAt)
		if err != nil {
			return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoot")
		}
		shoots = append(shoots, shoot)
	}

	if err := rows.Err(); err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "error during rows iteration")
	}

	for i := range shoots {
		clients, err := repo.getShootClients(ctx, shoots[i].Id)
		if err != nil {
			return nil, err
		}
		shoots[i].Clients = clients
	}

	return shoots, nil
}

func (repo *postgresShoot) getShootClients(ctx context.Context, shootID int) ([]model.ShootClientInfo, error) {
	query := `
SELECT
    cl.id, cl.first_name, cl.last_name, cl.phone, sc.is_main_client, sc.relationship_type
FROM shoot_clients sc
JOIN clients cl ON sc.client_id = cl.id
WHERE sc.shoot_id = $1
ORDER BY sc.is_main_client DESC, cl.first_name
`
	rows, err := repo.db.Pool.Query(ctx, query, shootID)
	defer rows.Close()

	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBQuery, "failed to get shoot clients")
	}

	clients := make([]model.ShootClientInfo, 0)
	for rows.Next() {
		var client model.ShootClientInfo
		err := rows.Scan(
			&client.ClientID,
			&client.FirstName,
			&client.LastName,
			&client.Phone,
			&client.IsMainClient,
			&client.RelationshipType,
		)
		if err != nil {
			return nil, myerrors.Wrap(err, myerrors.ErrCodeDBScan, "failed to scan shoot client")
		}
		clients = append(clients, client)
	}

	return clients, nil
}
