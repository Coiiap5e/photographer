package controllers

import (
	"context"

	"github.com/Coiiap5e/photographer/internal/api/controllers/request"
	"github.com/Coiiap5e/photographer/internal/api/controllers/response"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/validators"
)

type ClientController struct {
	clientService service.Client
}

func NewClientController(clientService service.Client) *ClientController {
	return &ClientController{
		clientService: clientService,
	}
}

func (cc *ClientController) CreateClient(ctx context.Context, req *request.CreateClientRequest) (*response.ClientResponse, error) {
	if err := validators.ValidateCreateClient(req); err != nil {
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

	return response.ToClientResponse(client), nil

}
