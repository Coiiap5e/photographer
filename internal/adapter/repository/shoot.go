package repository

import (
	"context"
	"errors"
	"time"

	myerrors "github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/infrastructure/database"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
	"github.com/jackc/pgx/v5"
)

type Shoot interface {
	AddShoot(ctx context.Context, shoot *model.Shoot, clients []*model.ShootClient) (*model.Shoot, error)
	DeleteShoot(ctx context.Context, id int) error
	GetShootByID(ctx context.Context, id int) (*model.Shoot, error)
	GetShoots(ctx context.Context) ([]model.Shoot, error)
	UpdateShoot(ctx context.Context, id int, shoot *model.Shoot, clients []*model.ShootClient) (*model.Shoot, error)
	UpdateShootDateTime(ctx context.Context, id int, patch *model.ShootDateTimePatch) (*model.Shoot, error)
	CountByDate(ctx context.Context, date time.Time) (int, error)
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

func (repo *postgresShoot) AddShoot(ctx context.Context, shoot *model.Shoot, clients []*model.ShootClient) (*model.Shoot, error) {
	tx, err := repo.db.Pool.Begin(ctx)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBTransaction, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	createdAt := repo.clock.Now()
	updatedAt := repo.clock.Now()

	query := `
INSERT INTO shoots
	(date, start_time, end_time, shoot_price, location, shoot_type, notes, created_at)
VALUES 
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
	id, created_at`

	shootToCreate := *shoot
	shootToCreate.CreatedAt = createdAt
	shootToCreate.UpdatedAt = updatedAt

	err = tx.QueryRow(ctx, query, shootToCreate.ShootDate,
		shootToCreate.StartTime, shootToCreate.EndTime, shootToCreate.ShootPrice, shootToCreate.ShootLocation,
		shootToCreate.ShootType, shootToCreate.Notes, shootToCreate.CreatedAt).Scan(&shootToCreate.Id, &shootToCreate.CreatedAt)

	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeShootCreate, "failed to create shoot")
	}

	err = repo.addAllClientsTx(ctx, tx, clients, shootToCreate.Id)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBTransaction, "failed to add clients")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBTransaction, "failed to commit transaction")
	}

	return &shootToCreate, nil
}

func (repo *postgresShoot) UpdateShoot(ctx context.Context, id int, shoot *model.Shoot, clients []*model.ShootClient) (*model.Shoot, error) {
	tx, err := repo.db.Pool.Begin(ctx)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBTransaction, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	updatedAt := repo.clock.Now()

	query := `
UPDATE shoots SET
	date = $1,
	start_time = $2,
	end_time = $3,
	shoot_price = $4,
	location = $5,
	shoot_type = $6,
	notes = $7,
	updated_at= $8
WHERE
    id = $9
RETURNING
	id, date, start_time, end_time, shoot_price, location, shoot_type, notes, created_at, updated_at`

	updatedShoot := *shoot
	updatedShoot.UpdatedAt = updatedAt

	err = tx.QueryRow(ctx, query,
		updatedShoot.ShootDate, updatedShoot.StartTime, updatedShoot.EndTime,
		updatedShoot.ShootPrice, updatedShoot.ShootLocation,
		updatedShoot.ShootType, updatedShoot.Notes, updatedShoot.UpdatedAt,
		id,
	).Scan(
		&updatedShoot.Id, &updatedShoot.ShootDate, &updatedShoot.StartTime,
		&updatedShoot.EndTime, &updatedShoot.ShootPrice, &updatedShoot.ShootLocation,
		&updatedShoot.ShootType, &updatedShoot.Notes,
		&updatedShoot.CreatedAt, &updatedShoot.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerrors.Wrap(err, myerrors.ErrCodeShootNotFound, "shoot not found")
		}
		return nil, myerrors.Wrap(err, myerrors.ErrCodeShootUpdate, "failed to update shoot")
	}

	_, err = tx.Exec(ctx, "DELETE FROM shoot_clients WHERE shoot_id = $1", updatedShoot.Id)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBTransaction, "failed to delete old clients")
	}

	err = repo.addAllClientsTx(ctx, tx, clients, updatedShoot.Id)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBTransaction, "failed to add clients")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBTransaction, "failed to commit transaction")
	}

	return &updatedShoot, nil
}

func (repo *postgresShoot) addAllClientsTx(ctx context.Context, tx pgx.Tx, clients []*model.ShootClient, shootID int) error {
	batch := &pgx.Batch{}

	for _, client := range clients {
		client.ShootID = shootID

		batch.Queue(`
INSERT INTO shoot_clients
    (shoot_id, client_id, is_main_client, relationship_type, created_at)
VALUES 
    ($1, $2, $3, $4, $5)
Returning created_at`,
			client.ShootID,
			client.ClientID,
			client.IsMainClient,
			client.RelationshipType,
			repo.clock.Now())
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < len(clients); i++ {
		err := results.QueryRow().Scan(&clients[i].CreatedAt)
		if err != nil {
			return myerrors.Wrap(err, myerrors.ErrCodeDBQuery, "failed to get created_at for shoot client")
		}
	}

	return results.Close()
}

func (repo *postgresShoot) UpdateShootDateTime(ctx context.Context, id int, patch *model.ShootDateTimePatch) (*model.Shoot, error) {
	tx, err := repo.db.Pool.Begin(ctx)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBTransaction, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	updatedAt := repo.clock.Now()
	query := `
UPDATE shoots SET
    date = $1,
    start_time = $2,
    end_time = $3,
    updated_at = $4
WHERE id = $5
RETURNING
	id, date, start_time, end_time, shoot_price, location, shoot_type, notes, created_at, updated_at
`
	var updatedShoot model.Shoot

	err = tx.QueryRow(ctx, query,
		patch.ShootDate,
		patch.StartTime,
		patch.EndTime,
		updatedAt,
		id,
	).Scan(
		&updatedShoot.Id,
		&updatedShoot.ShootDate,
		&updatedShoot.StartTime,
		&updatedShoot.EndTime,
		&updatedShoot.ShootPrice,
		&updatedShoot.ShootLocation,
		&updatedShoot.ShootType,
		&updatedShoot.Notes,
		&updatedShoot.CreatedAt,
		&updatedShoot.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerrors.Wrap(err, myerrors.ErrCodeShootNotFound, "shoot not found")
		}
		return nil, myerrors.Wrap(err, myerrors.ErrCodeShootUpdate, "failed to update shot")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBTransaction, "failed to commit transaction")
	}

	return &updatedShoot, nil
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
	shoot_price, location, shoot_type, notes, created_at
FROM shoots
WHERE id = $1`

	var shoot model.Shoot
	err := repo.db.Pool.QueryRow(ctx, query, id).Scan(
		&shoot.Id, &shoot.ShootDate, &shoot.StartTime,
		&shoot.EndTime, &shoot.ShootPrice, &shoot.ShootLocation,
		&shoot.ShootType, &shoot.Notes, &shoot.CreatedAt)

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
    s.id, s.date, s.start_time, s.end_time, 
    s.shoot_price, s.location, s.shoot_type, s.notes, 
    s.created_at, s.updated_at,
    sc.client_id, sc.is_main_client, sc.relationship_type,
    c.first_name, c.last_name, c.phone
FROM shoots s
INNER JOIN shoot_clients sc ON s.id = sc.shoot_id
INNER JOIN clients c ON sc.client_id = c.id
ORDER BY s.date DESC, s.created_at DESC
`
	rows, err := repo.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoots")
	}

	defer rows.Close()

	shootsMap := make(map[int]*model.Shoot)

	shoots := make([]model.Shoot, 0)

	for rows.Next() {
		var (
			shootID                         int
			shootDate                       time.Time
			startTime                       time.Time
			endTime                         time.Time
			shootPrice                      int
			shootLocation, shootType, notes string
			createdAt, updatedAt            time.Time
			clientID                        int
			isMainClient                    bool
			relationshipType                string
			firstName, lastName             string
			phone                           string
		)
		err := rows.Scan(
			&shootID, &shootDate, &startTime, &endTime,
			&shootPrice, &shootLocation, &shootType, &notes,
			&createdAt, &updatedAt,
			&clientID, &isMainClient, &relationshipType,
			&firstName, &lastName, &phone,
		)
		if err != nil {
			return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to scan shoot")
		}

		shoot, exists := shootsMap[shootID]
		if !exists {
			shoot = &model.Shoot{
				Id:            shootID,
				ShootDate:     shootDate,
				StartTime:     startTime,
				EndTime:       endTime,
				ShootPrice:    shootPrice,
				ShootLocation: shootLocation,
				ShootType:     shootType,
				Notes:         notes,
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
				Clients:       make([]model.ShootClientInfo, 0),
			}
			shootsMap[shootID] = shoot
		}

		client := model.ShootClientInfo{
			ClientID:         clientID,
			FirstName:        firstName,
			LastName:         lastName,
			Phone:            phone,
			IsMainClient:     isMainClient,
			RelationshipType: relationshipType,
		}
		shoot.Clients = append(shoot.Clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "error during rows iteration")
	}

	for _, shoot := range shootsMap {
		shoots = append(shoots, *shoot)
	}

	return shoots, nil
}

func (repo *postgresShoot) CountByDate(ctx context.Context, date time.Time) (int, error) {
	query := `
SELECT
	count(*)
FROM shoots
WHERE date = $1`

	var count int
	err := repo.db.Pool.QueryRow(ctx, query, date).Scan(&count)
	if err != nil {
		return 0, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoots count by date")
	}
	return count, nil
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
