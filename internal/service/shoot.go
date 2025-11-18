package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	myerrors "github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/repository"
	"github.com/Coiiap5e/photographer/internal/utils"
	"github.com/samber/lo"
)

type Shoot interface {
	CreateShoot(ctx context.Context, shoot *model.Shoot, clients []model.ShootClient) error
	DeleteShoot(ctx context.Context, id int) error
	GetShoots(ctx context.Context) error
	GetShootByID(ctx context.Context, id int) (*model.Shoot, error)
	GetShootsWithRelationshipType(ctx context.Context, relationshipType string) error
	GetShootsSortedByDate(ctx context.Context) error
}

type postgresShoot struct {
	shootRepo  repository.Shoot
	clientRepo repository.Client
	logger     *slog.Logger
}

func NewShoot(shootRepo repository.Shoot, clientRepo repository.Client, logger *slog.Logger) Shoot {
	return &postgresShoot{
		shootRepo:  shootRepo,
		clientRepo: clientRepo,
		logger:     logger,
	}
}

func (s *postgresShoot) CreateShoot(ctx context.Context, shoot *model.Shoot, clients []model.ShootClient) error {
	err := s.shootRepo.AddShoot(ctx, shoot, clients)
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

func (s *postgresShoot) GetShootsSortedByDate(ctx context.Context) error {
	allShoots, err := s.shootRepo.GetShoots(ctx)
	if err != nil {
		return err
	}

	sortedByDate := lo.GroupBy(allShoots, func(shoot model.Shoot) time.Time {
		return shoot.ShootDate
	})

	for date, shoots := range sortedByDate {
		fmt.Printf("Date: %s\n", date.Format("02.01.2006"))
		showShoots(shoots)
	}

	return nil
}

func (s *postgresShoot) GetShootsWithRelationshipType(ctx context.Context, relationshipType string) error {
	allShoots, err := s.shootRepo.GetShoots(ctx)
	if err != nil {
		return err
	}

	filteredShoots := lo.Filter(allShoots, func(shoot model.Shoot, _ int) bool {
		return lo.ContainsBy(shoot.Clients, func(client model.ShootClientInfo) bool {
			return client.RelationshipType == relationshipType
		})
	})

	showShoots(filteredShoots)

	return nil
}

func showShoots(shoots []model.Shoot) {
	if len(shoots) == 0 {
		fmt.Println("No shoots found")
		return
	}

	fmt.Printf("%-3s %-10s %-8s %-8s %-6s %-25s %-10s %-25s %-10s\n",
		"ID", "Date", "Start", "End", "Price", "Location",
		"Type", "Notes", "Created")

	fmt.Println(strings.Repeat("-", 148))

	for _, shoot := range shoots {
		fmt.Printf("%-3d %-10s %-8s %-8s %-6d %-25s %-10s %-25s %-10s\n",
			shoot.Id,
			shoot.ShootDate.Format("02.01.2006"),
			shoot.StartTime.Format("15:04"),
			shoot.EndTime.Format("15:04"),
			shoot.ShootPrice,
			shoot.ShootLocation,
			shoot.ShootType,
			shoot.Notes,
			shoot.CreatedAt,
		)

		if len(shoot.Clients) > 0 {
			fmt.Println("Clients:")
			for _, client := range shoot.Clients {
				mainIndicator := ""
				if client.IsMainClient {
					mainIndicator = "MAIN"
				}
				fmt.Printf("     * %s %s (%s) - %s (%s)\n",
					client.FirstName,
					client.LastName,
					client.Phone,
					client.RelationshipType,
					mainIndicator,
				)
			}
			fmt.Println(strings.Repeat("-", 148))
		} else {
			fmt.Println("No clients")
		}
	}

}
