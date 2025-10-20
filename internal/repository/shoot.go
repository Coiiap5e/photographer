package repository

import (
	"context"
	"errors"

	"github.com/Coiiap5e/photographer/internal/database"
	myerrors "github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/jackc/pgx/v5"
)

type Shoot interface {
	AddShoot(ctx context.Context, shoot *model.Shoot) error
	DeleteShoot(ctx context.Context, id int) error
	GetShootByID(ctx context.Context, id int) (*model.Shoot, error)
	GetShoots(ctx context.Context) ([]model.Shoot, error)
}

type postgresShoot struct {
	db *database.DB
}

func NewShoot(db *database.DB) Shoot {
	return &postgresShoot{db: db}
}

func (repo *postgresShoot) AddShoot(ctx context.Context, shoot *model.Shoot) error {
	query := `
INSERT INTO shoots
	(client_id, date, start_time, end_time, shoot_price, location, client_first_name, client_last_name, shoot_type, notes)
VALUES 
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING
	id, created_at`

	err := repo.db.Pool.QueryRow(ctx, query, shoot.ClientId, shoot.ShootDate,
		shoot.StartTime, shoot.EndTime, shoot.ShootPrice, shoot.ShootLocation,
		shoot.ClientFirstName, shoot.ClientLastName, shoot.ShootType, shoot.Notes).Scan(&shoot.Id, &shoot.CreatedAt)

	if err != nil {
		return myerrors.Wrap(err, myerrors.ErrCodeShootCreate, "failed to create client")
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
	id, client_id, date, start_time, end_time, 
	shoot_price, location, client_first_name, 
	client_last_name, shoot_type, notes
FROM shoots
WHERE id = $1`

	var shoot model.Shoot
	err := repo.db.Pool.QueryRow(ctx, query, id).Scan(
		&shoot.Id, &shoot.ClientId, &shoot.ShootDate, &shoot.StartTime,
		&shoot.EndTime, &shoot.ShootPrice, &shoot.ShootLocation,
		&shoot.ClientFirstName, &shoot.ClientLastName, &shoot.ShootType,
		&shoot.Notes)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerrors.New(myerrors.ErrCodeShootNotFound, "shoot not found")
		}
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoot")
	}

	return &shoot, nil
}

func (repo *postgresShoot) GetShoots(ctx context.Context) ([]model.Shoot, error) {
	query := `
SELECT 
    id, client_id, date, start_time, end_time, 
    shoot_price, location, client_first_name, 
    client_last_name, shoot_type, notes, created_at
FROM shoots`

	rows, err := repo.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoots")
	}

	defer rows.Close()

	shoots := make([]model.Shoot, 0)

	for rows.Next() {
		var shoot model.Shoot
		err := rows.Scan(
			&shoot.Id, &shoot.ClientId, &shoot.ShootDate, &shoot.StartTime,
			&shoot.EndTime, &shoot.ShootPrice, &shoot.ShootLocation,
			&shoot.ClientFirstName, &shoot.ClientLastName, &shoot.ShootType,
			&shoot.Notes, &shoot.CreatedAt)
		if err != nil {
			return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoot")
		}
		shoots = append(shoots, shoot)
	}

	if err := rows.Err(); err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "error during rows iteration")
	}

	return shoots, nil
}
