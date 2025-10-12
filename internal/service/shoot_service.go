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

type ShootService struct {
	shootRepo *repository.ShootRepository
}

func NewShootService(shootRepo *repository.ShootRepository) *ShootService {
	return &ShootService{shootRepo: shootRepo}
}

func (s *ShootService) CreateShoot(ctx context.Context) error {
	shootDate, startTime, endTime := utils.InputShootDate()

	shoot := &model.Shoot{
		ClientId:        utils.InputId("Client_id"),
		ShootDate:       shootDate,
		StartTime:       startTime,
		EndTime:         endTime,
		ShootPrice:      utils.InputInt("Shoot price"),
		ShootLocation:   utils.InputStringRequired("Location"),
		ClientFirstName: utils.InputStringRequired("Client first name"),
		ClientLastName:  utils.InputStringRequired("Client last name"),
		ShootType:       utils.InputStringRequired("Shoot type"),
		Notes:           utils.InputString("Notes"),
	}

	err := s.shootRepo.AddShoot(ctx, shoot)
	if err != nil {
		return err
	}

	log.Printf("Shoot added successfully")
	return nil
}

func (s *ShootService) DeleteShoot(ctx context.Context) error {
	id := utils.InputId("ID of the shoot")

	shoot, err := s.shootRepo.GetShootByID(ctx, id)
	if err != nil {
		return err
	}

	fmt.Println("Confirm the deletion")
	fmt.Printf("Deleting shoot: %s start: %s end: %s\n",
		shoot.ShootDate.Format("02.01.2006"),
		shoot.StartTime.Format("15:04"),
		shoot.EndTime.Format("15:04"))
	fmt.Printf("Location: %s. ShootType: %s. Price: %d\n",
		shoot.ShootLocation, shoot.ShootType, shoot.ShootPrice)
	fmt.Printf("Client id: %d name: %s %s\n", shoot.ClientId,
		shoot.ClientFirstName, shoot.ClientLastName)
	if shoot.Notes != "" {
		fmt.Printf("Notes: %s\n", shoot.Notes)
	}

	confirm := utils.InputStringRequired("Are you sure you want to delete the shoot? (y/n)")
	if confirm == "n" || confirm == "N" {
		return errors.New(errors.ErrCodeValidation, "deletion cancelled")
	}

	err = s.shootRepo.DeleteShoot(ctx, id)
	if err != nil {
		return err
	}

	log.Printf("Shoot %s with %s %s deleted successfully \n",
		shoot.StartTime.Format("02.01.2006 15:04"),
		shoot.ClientFirstName, shoot.ClientLastName)
	return nil
}

func (s *ShootService) GetShoots(ctx context.Context) error {
	shoots, err := s.shootRepo.GetShoots(ctx)
	if err != nil {
		return err
	}

	s.showShoots(shoots)

	return nil
}

func (s *ShootService) showShoots(shoots []model.Shoot) {
	if len(shoots) == 0 {
		fmt.Println("No shoots found")
		return
	}

	fmt.Printf("%-3s %-9s %-10s %-8s %-8s %-6s %-12s %-12s %-12s %-10s %-15s %-10s\n",
		"ID", "Client ID", "Date", "Start", "End", "Price", "Location",
		"First name", "Last name", "Type", "Notes", "Created")

	fmt.Println(strings.Repeat("-", 125))

	for _, shoot := range shoots {
		fmt.Printf("%-3d %-9d %-10s %-8s %-8s %-6d %-12s %-12s %-12s %-10s %-15s %-10s\n",
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
