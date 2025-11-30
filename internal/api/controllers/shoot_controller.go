package controllers

import (
	"context"
	"fmt"

	"github.com/Coiiap5e/photographer/internal/api/controllers/request"
	"github.com/Coiiap5e/photographer/internal/api/controllers/response"
	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
	"github.com/Coiiap5e/photographer/internal/validators"
)

type ShootController struct {
	shootService  service.Shoot
	clientService service.Client
}

func NewShootController(
	shootService service.Shoot,
	clientService service.Client,
) *ShootController {
	return &ShootController{
		shootService:  shootService,
		clientService: clientService,
	}
}

func (sc *ShootController) CreateShoot(ctx context.Context, req *request.CreateShootRequest, clock *clock.Clock) (*response.ShootResponse, error) {
	if err := validators.ValidateCreateShoot(req, clock); err != nil {
		return nil, err
	}

	for _, clientInfo := range req.Clients {
		_, err := sc.clientService.GetClientByID(ctx, clientInfo.ClientID)
		if err != nil {
			if errors.IsErrorCode(err, errors.ErrCodeClientNotFound) {
				return nil, errors.New(errors.ErrCodeClientNotFound,
					fmt.Sprintf("client with ID %d not found", clientInfo.ClientID))
			}
			return nil, errors.Wrap(err, errors.ErrCodeDBSelect,
				fmt.Sprintf("failed to get client with ID %d", clientInfo.ClientID))
		}
	}

	clients := request.ToShootClientDomain(req.Clients)

	shoot := &model.Shoot{
		ShootDate:     req.ShootDate,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		ShootPrice:    req.ShootPrice,
		ShootLocation: req.ShootLocation,
		ShootType:     req.ShootType,
		Notes:         req.Notes,
	}

	err := sc.shootService.CreateShoot(ctx, shoot, clients)
	if err != nil {
		return nil, err
	}

	// TODO: мы не знаем id съемки, надо изменить логику создания съемки в сервисе и репозитории (возвращать или id или все данные съемки)
	return response.ToShootResponse(shoot), nil
}
