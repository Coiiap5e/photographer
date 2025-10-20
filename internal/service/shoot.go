package service

import (
	"context"
	"fmt"
	"strings"

	myerrors "github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/repository"
	"github.com/Coiiap5e/photographer/internal/utils"
)

type Shoot interface {
	CreateShoot(ctx context.Context, shoot *model.Shoot) error
	DeleteShoot(ctx context.Context, id int) error
	GetShoots(ctx context.Context) error
	GetShootByID(ctx context.Context, id int) (*model.Shoot, error)
}

type postgresShoot struct {
	shootRepo  repository.Shoot
	clientRepo repository.Client
}

func NewShoot(shootRepo repository.Shoot, clientRepo repository.Client) Shoot {
	return &postgresShoot{
		shootRepo:  shootRepo,
		clientRepo: clientRepo,
	}
}

func (s *postgresShoot) CreateShoot(ctx context.Context, shoot *model.Shoot) error {
	err := s.shootRepo.AddShoot(ctx, shoot)
	if err != nil {
		return err
	}

	return nil
}

func (s *postgresShoot) GetShootByID(ctx context.Context, id int) (*model.Shoot, error) {
	shoot, err := s.shootRepo.GetShootByID(ctx, id)
	if err != nil {
		if myerrors.IsErrorCode(err, myerrors.ErrCodeShootNotFound) {
			return nil, myerrors.Wrap(err, myerrors.ErrCodeShootNotFound, "shoot not found")
		}
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoot")
	}

	return shoot, nil
}

func (s *postgresShoot) DeleteShoot(ctx context.Context, id int) error {
	for {
		confirm := utils.InputStringRequired("Are you sure you want to delete the shoot? (y/n)")
		if confirm == "n" || confirm == "N" {
			return myerrors.New(myerrors.ErrCodeValidation, "deletion cancelled")
		} else if confirm == "y" || confirm == "Y" {
			break
		} else {
			fmt.Println("Press wrong button: enter (y/n)")
		}
	}

	err := s.shootRepo.DeleteShoot(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *postgresShoot) GetShoots(ctx context.Context) error {
	shoots, err := s.shootRepo.GetShoots(ctx)
	if err != nil {
		return err
	}

	showShoots(shoots)

	return nil
}

func showShoots(shoots []model.Shoot) {
	if len(shoots) == 0 {
		fmt.Println("No shoots found")
		return
	}

	fmt.Printf("%-3s %-9s %-10s %-8s %-8s %-6s %-25s %-12s %-12s %-10s %-25s %-10s\n",
		"ID", "Client ID", "Date", "Start", "End", "Price", "Location",
		"First name", "Last name", "Type", "Notes", "Created")

	fmt.Println(strings.Repeat("-", 148))

	for _, shoot := range shoots {
		fmt.Printf("%-3d %-9d %-10s %-8s %-8s %-6d %-25s %-12s %-12s %-10s %-25s %-10s\n",
			shoot.Id,
			shoot.ClientId,
			shoot.ShootDate.Format("02.01.2006"),
			shoot.StartTime.Format("15:04"),
			shoot.EndTime.Format("15:04"),
			shoot.ShootPrice,
			shoot.ShootLocation,
			shoot.ClientFirstName,
			shoot.ClientLastName,
			shoot.ShootType,
			shoot.Notes,
			shoot.CreatedAt,
		)
	}
}
