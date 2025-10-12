package repository

import (
	"context"

	"github.com/Coiiap5e/photographer/internal/database"
	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/model"
)

type ShootRepository struct {
	db *database.DB
}

func NewShootRepository(db *database.DB) *ShootRepository {
	return &ShootRepository{db: db}
}

func (repo *ShootRepository) AddShoot(ctx context.Context, shoot *model.Shoot) error {
	query := `
INSERT INTO shoots
	(client_id, date, start_time, end_time, shoot_price, location, client_first_name, client_last_name, shoot_type, notes)
VALUES 
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
	id, created_at`

	err := repo.db.Pool.QueryRow(ctx, query, shoot.ClientId, shoot.ShootDate,
		shoot.StartTime, shoot.EndTime, shoot.ShootPrice, shoot.ShootLocation,
		shoot.ClientFirstName, shoot.ClientLastName, shoot.ShootType, shoot.Notes).Scan(&shoot.Id, &shoot.CreatedAt)

	if err != nil {
		return errors.Wrap(err, errors.ErrCodeShootCreate, "failed to create client")
	}

	return nil
}
