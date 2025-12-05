package controllers

import (
	"context"

	"github.com/Coiiap5e/photographer/internal/api/controllers/request"
	"github.com/Coiiap5e/photographer/internal/api/controllers/response"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
	"github.com/Coiiap5e/photographer/internal/validators"
)

type ShootController struct {
	shootService  service.Shoot
	clientService service.Client
	clock         *clock.Clock
}

func NewShootController(
	shootService service.Shoot,
	clientService service.Client,
	clock *clock.Clock,
) *ShootController {
	return &ShootController{
		shootService:  shootService,
		clientService: clientService,
		clock:         clock,
	}
}

func (sc *ShootController) CreateShoot(ctx context.Context, req *request.CreateShootRequest) (*response.ShootResponse, error) {
	if err := validators.ValidateCreateShoot(req, sc.clock); err != nil {
		return nil, err
	}

	clients := request.ToShootClientDomain(req.Clients)

	newShoot := &model.Shoot{
		ShootDate:     req.ShootDate,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		ShootPrice:    req.ShootPrice,
		ShootLocation: req.ShootLocation,
		ShootType:     req.ShootType,
		Notes:         req.Notes,
	}

	createdShoot, err := sc.shootService.CreateShoot(ctx, newShoot, clients)
	if err != nil {
		return nil, err
	}

	return response.ToShootResponse(createdShoot), nil
}

func (sc *ShootController) UpdateShoot(ctx context.Context, id int, req *request.CreateShootRequest) (*response.ShootResponse, error) {
	if err := validators.ValidateCreateShoot(req, sc.clock); err != nil {
		return nil, err
	}

	if err := validators.ValidateID(id); err != nil {
		return nil, err
	}

	clients := request.ToShootClientDomain(req.Clients)

	newShoot := &model.Shoot{
		ShootDate:     req.ShootDate,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		ShootPrice:    req.ShootPrice,
		ShootLocation: req.ShootLocation,
		ShootType:     req.ShootType,
		Notes:         req.Notes,
	}

	updatedShoot, err := sc.shootService.UpdateShoot(ctx, id, newShoot, clients)
	if err != nil {
		return nil, err
	}

	return response.ToShootResponse(updatedShoot), nil
}

func (sc *ShootController) UpdateShootDateTime(ctx context.Context, id int, req *request.UpdateShootDateTimeRequest) (*response.ShootResponse, error) {
	if err := validators.ValidateDate(req, sc.clock); err != nil {
		return nil, err
	}

	if err := validators.ValidateID(id); err != nil {
		return nil, err
	}

	newUpdatedDateTime := &model.ShootDateTimePatch{
		ShootDate: req.ShootDate,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	updatedShoot, err := sc.shootService.UpdateShootDateTime(ctx, id, newUpdatedDateTime)
	if err != nil {
		return nil, err
	}

	return response.ToShootResponse(updatedShoot), nil

}

func (sc *ShootController) GetShootByID(ctx context.Context, id int) (*response.ShootResponse, error) {
	if err := validators.ValidateID(id); err != nil {
		return nil, err
	}

	shoot, err := sc.shootService.GetShootByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return response.ToShootResponse(shoot), nil

}

func (sc *ShootController) DeleteShoot(ctx context.Context, id int) error {
	if err := validators.ValidateID(id); err != nil {
		return err
	}

	err := sc.shootService.DeleteShoot(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
