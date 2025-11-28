package controllers

import (
	"context"
	"fmt"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/validators"
)

type ShootController struct {
	shootService  service.Shoot
	clientService service.Client
	validator     *validators.ShootValidator
}

func NewShootController(
	shootService service.Shoot,
	clientService service.Client,
	validator *validators.ShootValidator,
) *ShootController {
	return &ShootController{
		shootService:  shootService,
		clientService: clientService,
		validator:     validator,
	}
}

func (sc *ShootController) CreateShoot(ctx context.Context, req *model.CreateShootRequest) (*model.ShootResponse, error) {
	if err := sc.validator.ValidateCreateShoot(req); err != nil {
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

	clients := sc.toShootClientDomain(req.Clients)

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
	return sc.toShootResponse(shoot), nil
}

func (sc *ShootController) toShootClientDomain(clientRequests []model.ShootClientRequest) []*model.ShootClient {
	clients := make([]*model.ShootClient, len(clientRequests))
	for i, clientRequest := range clientRequests {
		clients[i] = &model.ShootClient{
			ClientID:         clientRequest.ClientID,
			IsMainClient:     clientRequest.IsMainClient,
			RelationshipType: clientRequest.RelationshipType,
		}
	}
	return clients
}

func (sc *ShootController) toShootClientInfoResponse(clients []model.ShootClientInfo) []model.ShootClientInfoResponse {
	clientsResponse := make([]model.ShootClientInfoResponse, len(clients))
	for i, clientInfo := range clients {
		clientsResponse[i] = model.ShootClientInfoResponse{
			ClientID:         clientInfo.ClientID,
			FirstName:        clientInfo.FirstName,
			LastName:         clientInfo.LastName,
			Phone:            clientInfo.Phone,
			IsMainClient:     clientInfo.IsMainClient,
			RelationshipType: clientInfo.RelationshipType,
		}
	}

	return clientsResponse
}

func (sc *ShootController) toShootResponse(shoot *model.Shoot) *model.ShootResponse {
	clients := sc.toShootClientInfoResponse(shoot.Clients)
	return &model.ShootResponse{
		ID:            shoot.Id,
		ShootDate:     shoot.ShootDate,
		StartTime:     shoot.StartTime,
		EndTime:       shoot.EndTime,
		ShootPrice:    shoot.ShootPrice,
		ShootLocation: shoot.ShootLocation,
		ShootType:     shoot.ShootType,
		Notes:         shoot.Notes,
		CreatedAt:     shoot.CreatedAt,
		Clients:       clients,
	}
}
