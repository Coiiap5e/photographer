package handlers

import (
	"net/http"

	"github.com/Coiiap5e/photographer/internal/api/controllers"
	"github.com/Coiiap5e/photographer/internal/api/controllers/request"
	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	clientController *controllers.ClientController
}

func NewClientHandler(clientController *controllers.ClientController) *ClientHandler {
	return &ClientHandler{
		clientController: clientController,
	}
}

func (h *ClientHandler) CreateClient(c *gin.Context) {
	var req request.CreateClientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}

	response, err := h.clientController.CreateClient(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.IsErrorCode(err, errors.ErrCodeValidation):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "VALIDATION_ERROR",
				"message": err.Error(),
			})
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

	c.JSON(http.StatusCreated, response)
}
