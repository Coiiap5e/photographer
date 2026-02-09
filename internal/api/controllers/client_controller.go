package controllers

import (
	"net/http"

	"github.com/Coiiap5e/photographer/internal/api/controllers/dto"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/validators"
	"github.com/gin-gonic/gin"
)

type ClientController struct {
	clientService service.Client
}

func NewClientController(clientService service.Client) *ClientController {
	return &ClientController{
		clientService: clientService,
	}
}

func (cc *ClientController) CreateClient(c *gin.Context) {
	var req dto.CreateClientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := validators.ValidateCreateClient(&req); err != nil {
		_ = c.Error(err)
		return
	}

	client := &model.Client{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Phone:            req.Phone,
		SocialNetworkUrl: req.SocialNetworkUrl,
	}

	ctx := c.Request.Context()
	createdClient, err := cc.clientService.CreateClient(ctx, client)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToClientResponse(createdClient))
}
