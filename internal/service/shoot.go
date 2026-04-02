package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Coiiap5e/photographer/internal/adapter/currency"
	"github.com/Coiiap5e/photographer/internal/adapter/repository"
	myerrors "github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
	"github.com/samber/lo"
)

type Shoot interface {
	CreateShoot(ctx context.Context, shoot *model.Shoot, clients []*model.ShootClient) (*model.Shoot, error)
	DeleteShoot(ctx context.Context, id int) error
	GetShoots(ctx context.Context) ([]model.Shoot, error)
	GetShootByID(ctx context.Context, id int) (*model.Shoot, error)
	GetShootsWithRelationshipType(ctx context.Context, relationshipType string) error
	GetShootsSortedByDate(ctx context.Context) error
	UpdateShoot(ctx context.Context, id int, shoot *model.Shoot, clients []*model.ShootClient) (*model.Shoot, error)
	UpdateShootDateTime(ctx context.Context, id int, patch *model.ShootDateTimePatch) (*model.Shoot, error)
	GetShootsForNextTwoDays(ctx context.Context) ([]model.Shoot, error)
	GetShootsForNextTwoHours(ctx context.Context) ([]model.Shoot, error)
}

type postgresShoot struct {
	shootRepo       repository.Shoot
	clientService   Client
	currencyService currency.Service
	logger          *slog.Logger
	clock           *clock.Clock
}

func NewShoot(shootRepo repository.Shoot, clientService Client, currencyService currency.Service, logger *slog.Logger, clock *clock.Clock) Shoot {
	return &postgresShoot{
		shootRepo:       shootRepo,
		clientService:   clientService,
		currencyService: currencyService,
		logger:          logger,
		clock:           clock,
	}
}

func (s *postgresShoot) CreateShoot(ctx context.Context, shoot *model.Shoot, clients []*model.ShootClient) (*model.Shoot, error) {
	for _, client := range clients {
		_, err := s.clientService.GetClientByID(ctx, client.ClientID)
		if err != nil {
			if myerrors.IsErrorCode(err, myerrors.ErrCodeClientNotFound) {
				return nil, myerrors.New(myerrors.ErrCodeClientNotFound,
					fmt.Sprintf("client with ID %d not found", client.ClientID))
			}
			return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect,
				fmt.Sprintf("failed to get client with ID %d", client.ClientID))
		}
	}

	createdShoot, err := s.shootRepo.AddShoot(ctx, shoot, clients)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeShootCreate, "failed to create shoot")
	}

	return createdShoot, nil
}

func (s *postgresShoot) UpdateShoot(ctx context.Context, id int, shoot *model.Shoot, clients []*model.ShootClient) (*model.Shoot, error) {
	_, err := s.shootRepo.GetShootByID(ctx, id)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeShootNotFound, "failed to get shoot for update")
	}

	for _, client := range clients {
		_, err = s.clientService.GetClientByID(ctx, client.ClientID)
		if err != nil {
			if myerrors.IsErrorCode(err, myerrors.ErrCodeClientNotFound) {
				return nil, myerrors.New(myerrors.ErrCodeClientNotFound,
					fmt.Sprintf("client with ID %d not found", client.ClientID))
			}
			return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect,
				fmt.Sprintf("failed to get client with ID %d", client.ClientID))
		}
	}

	updatedShoot, err := s.shootRepo.UpdateShoot(ctx, id, shoot, clients)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeShootUpdate, "failed to update shoot")
	}

	return updatedShoot, nil
}

func (s *postgresShoot) UpdateShootDateTime(ctx context.Context, id int, patch *model.ShootDateTimePatch) (*model.Shoot, error) {
	_, err := s.shootRepo.GetShootByID(ctx, id)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeShootNotFound, "failed to get shoot for date/time update")
	}

	updatedShoot, err := s.shootRepo.UpdateShootDateTime(ctx, id, patch)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeShootUpdate, "failed to update shoot date/time")
	}

	return updatedShoot, nil
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
	_, err := s.GetShootByID(ctx, id)

	if err != nil {
		if myerrors.IsErrorCode(err, myerrors.ErrCodeShootNotFound) {
			return myerrors.Wrap(err, myerrors.ErrCodeShootNotFound, "shoot not found")
		}
		return myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoot")
	}

	err = s.shootRepo.DeleteShoot(ctx, id)
	if err != nil {
		return myerrors.Wrap(err, myerrors.ErrCodeShootDelete, "failed to delete shoot")
	}

	return nil
}

func (s *postgresShoot) GetShoots(ctx context.Context) ([]model.Shoot, error) {
	shoots, err := s.shootRepo.GetShoots(ctx)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoots")
	}

	return shoots, nil
}

func (s *postgresShoot) GetShootsSortedByDate(ctx context.Context) error {
	allShoots, err := s.shootRepo.GetShoots(ctx)
	if err != nil {
		return myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoots")
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
		return myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoots")
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

func (s *postgresShoot) GetShootsForNextTwoDays(ctx context.Context) ([]model.Shoot, error) {
	now := s.clock.Now()
	twoDaysLater := now.Add(48 * time.Hour)
	shoots, err := s.shootRepo.GetShootsByTimeRange(ctx, now, twoDaysLater)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoots for next two days")
	}
	return s.enrichShootsWithUSDPrice(ctx, shoots)
}

func (s *postgresShoot) GetShootsForNextTwoHours(ctx context.Context) ([]model.Shoot, error) {
	now := s.clock.Now()
	twoHoursLater := now.Add(2 * time.Hour)
	shoots, err := s.shootRepo.GetShootsByTimeRange(ctx, now, twoHoursLater)
	if err != nil {
		return nil, myerrors.Wrap(err, myerrors.ErrCodeDBSelect, "failed to get shoots for next two hours")
	}
	return s.enrichShootsWithUSDPrice(ctx, shoots)
}

func (s *postgresShoot) enrichShootsWithUSDPrice(ctx context.Context, shoots []model.Shoot) ([]model.Shoot, error) {
	if len(shoots) == 0 {
		return shoots, nil
	}

	rate, err := s.currencyService.GetUSDRate(ctx)
	if err != nil {
		s.logger.Error("failed to get USD rate, returning shoots without USD price", "error", err)
		return shoots, nil
	}

	if rate == 0 {
		s.logger.Error("got a zero USD rate, returning shoots without USD price")
		return shoots, nil
	}

	for i := range shoots {
		shoots[i].PriceUSD = float64(shoots[i].ShootPrice) / rate
	}

	return shoots, nil
}
