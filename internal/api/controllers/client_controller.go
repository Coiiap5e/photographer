package controllers

import (
	"net/http"

	"github.com/Coiiap5e/photographer/internal/api/controllers/request"
	"github.com/Coiiap5e/photographer/internal/errors"
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
	var req request.CreateClientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}

	if err := validators.ValidateCreateClient(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "VALIDATION_ERROR",
			"message": err.Error(),
		})
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
		switch {
		case errors.IsErrorCode(err, errors.ErrCodeClientCreate):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   "CREATE_CLIENT_ERROR",
				"message": "Failed to create client",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_ERROR",
				"message": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, request.ToClientResponse(createdClient))
}
