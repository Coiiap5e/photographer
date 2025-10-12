package service

import (
	"context"
	"log"

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
