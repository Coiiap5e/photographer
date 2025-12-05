package handlers

import (
	"net/http"

	"github.com/Coiiap5e/photographer/internal/api/controllers"
	"github.com/Coiiap5e/photographer/internal/api/controllers/request"
	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/gin-gonic/gin"
)

type ShootHandler struct {
	shootController *controllers.ShootController
}

func NewShootHandler(shootController *controllers.ShootController) *ShootHandler {
	return &ShootHandler{
		shootController: shootController,
	}
}

func (h *ShootHandler) CreateShoot(c *gin.Context) {
	var req request.CreateShootRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}

	response, err := h.shootController.CreateShoot(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.IsErrorCode(err, errors.ErrCodeValidation):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "VALIDATION_ERROR",
				"message": err.Error(),
			})
		case errors.IsErrorCode(err, errors.ErrCodeClientNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "CLIENT_NOT_FOUND",
				"message": err.Error(),
			})
		case errors.IsErrorCode(err, errors.ErrCodeShootCreate):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   "CREATE_SHOOT_ERROR",
				"message": "Failed to create shoot",
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

func (h *ShootHandler) UpdateShoot(c *gin.Context) {
	var uriParams struct {
		ID int `uri:"id" binding:"min=1"`
	}

	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid shoot ID in URL",
		})
		return
	}

	var req request.CreateShootRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}

	response, err := h.shootController.UpdateShoot(c.Request.Context(), uriParams.ID, &req)
	if err != nil {
		switch {
		case errors.IsErrorCode(err, errors.ErrCodeValidation):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "VALIDATION_ERROR",
				"message": err.Error(),
			})
		case errors.IsErrorCode(err, errors.ErrCodeClientNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "CLIENT_NOT_FOUND",
				"message": err.Error(),
			})
		case errors.IsErrorCode(err, errors.ErrCodeShootUpdate):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   "UPDATE_SHOOT_ERROR",
				"message": "Failed to update shoot",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_ERROR",
				"message": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ShootHandler) UpdateShootDateTime(c *gin.Context) {
	var uriParams struct {
		ID int `uri:"id" binding:"min=1"`
	}

	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid shoot ID in URL",
		})
		return
	}

	var req request.UpdateShootDateTimeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}

	response, err := h.shootController.UpdateShootDateTime(c.Request.Context(), uriParams.ID, &req)
	if err != nil {
		switch {
		case errors.IsErrorCode(err, errors.ErrCodeValidation):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "VALIDATION_ERROR",
				"message": err.Error(),
			})
		case errors.IsErrorCode(err, errors.ErrCodeShootUpdate):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   "UPDATE_SHOOT_ERROR",
				"message": "Failed to update shoot",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "INTERNAL_ERROR",
				"message": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ShootHandler) GetShootByID(c *gin.Context) {
	var uriParams request.GetShootByIDRequest

	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid shoot ID",
		})
		return
	}

	response, err := h.shootController.GetShootByID(c.Request.Context(), uriParams.ID)

	if err != nil {
		if errors.IsErrorCode(err, errors.ErrCodeShootNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "SHOT_NOT_FOUND",
				"message": "Shoot not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to get shoot",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *ShootHandler) GetShoots(c *gin.Context) {

	response, err := h.shootController.GetShoots(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to get shoot",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ShootHandler) DeleteShoot(c *gin.Context) {
	var uriParams struct {
		ID int `uri:"id" binding:"min=1"`
	}

	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid shoot ID",
		})
		return
	}

	err := h.shootController.DeleteShoot(c.Request.Context(), uriParams.ID)
	if err != nil {
		if errors.IsErrorCode(err, errors.ErrCodeShootNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "SHOOT_NOT_FOUND",
				"message": "Shoot not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "DELETE_SHOOT_ERROR",
			"message": "Failed to delete shoot",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shoot deleted successfully",
		"deleted": true,
	})

}
