package controllers

import (
	"context"

	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/validators"
)

type ClientController struct {
	clientService service.Client
	validator     *validators.ClientValidator
}

func NewClientController(clientService service.Client, validator *validators.ClientValidator) *ClientController {
	return &ClientController{
		clientService: clientService,
		validator:     validator,
	}
}

func (cc *ClientController) CreateClient(ctx context.Context, req *model.CreateClientRequest) (*model.ClientResponse, error) {
	if err := cc.validator.ValidateCreateClient(req); err != nil {
		return nil, err
	}

	client := &model.Client{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Phone:            req.Phone,
		SocialNetworkUrl: req.SocialNetworkUrl,
	}

	err := cc.clientService.CreateClient(ctx, client)
	if err != nil {
		return nil, err
	}

	return cc.toClientResponse(client), nil

}

func (cc *ClientController) toClientResponse(client *model.Client) *model.ClientResponse {
	return &model.ClientResponse{
		ID:               client.Id,
		FirstName:        client.FirstName,
		LastName:         client.LastName,
		Phone:            client.Phone,
		SocialNetworkUrl: client.SocialNetworkUrl,
		CreatedAt:        client.CreatedAt,
	}
}

func (cc *ClientController) toClientsResponse(clients []*model.Client) []model.ClientResponse {
	responses := make([]model.ClientResponse, len(clients))
	for i, client := range clients {
		responses[i] = *cc.toClientResponse(client)
	}
	return responses
}
