package controllers

import (
	"net/http"

	"github.com/Coiiap5e/photographer/internal/api/controllers/request"
	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/service"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
	"github.com/Coiiap5e/photographer/internal/validators"
	"github.com/gin-gonic/gin"
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

func (sc *ShootController) CreateShoot(c *gin.Context) {
	var req request.CreateShootRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}
	if err := validators.ValidateCreateShoot(&req, sc.clock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
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

	createdShoot, err := sc.shootService.CreateShoot(c.Request.Context(), newShoot, clients)
	if err != nil {
		switch {
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

	c.JSON(http.StatusCreated, request.ToShootResponse(createdShoot))
}

func (sc *ShootController) UpdateShoot(c *gin.Context) {
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
	if err := validators.ValidateCreateShoot(&req, sc.clock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
	}
	if err := validators.ValidateID(uriParams.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
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

	updatedShoot, err := sc.shootService.UpdateShoot(c.Request.Context(), uriParams.ID, newShoot, clients)
	if err != nil {
		switch {
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

	c.JSON(http.StatusOK, request.ToShootResponse(updatedShoot))
}

func (sc *ShootController) UpdateShootDateTime(c *gin.Context) {
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
	if err := validators.ValidateDate(&req, sc.clock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
	}
	if err := validators.ValidateID(uriParams.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
	}

	newUpdatedDateTime := &model.ShootDateTimePatch{
		ShootDate: req.ShootDate,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	updatedShoot, err := sc.shootService.UpdateShootDateTime(c.Request.Context(), uriParams.ID, newUpdatedDateTime)
	if err != nil {
		switch {
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

	c.JSON(http.StatusOK, request.ToShootResponse(updatedShoot))
}

func (sc *ShootController) GetShootByID(c *gin.Context) {
	var uriParams request.GetShootByIDRequest
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_REQUEST",
			"message": "Invalid shoot ID",
		})
		return
	}

	if err := validators.ValidateID(uriParams.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
	}

	shoot, err := sc.shootService.GetShootByID(c.Request.Context(), uriParams.ID)
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
	c.JSON(http.StatusOK, request.ToShootResponse(shoot))
}

func (sc *ShootController) GetShoots(c *gin.Context) {
	shoots, err := sc.shootService.GetShoots(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to get shoot",
		})
		return
	}

	responseShoots := make([]*request.ShootResponse, 0)
	for _, shoot := range shoots {
		responseShoots = append(responseShoots, request.ToShootResponse(&shoot))
	}

	c.JSON(http.StatusOK, responseShoots)
}

func (sc *ShootController) DeleteShoot(c *gin.Context) {
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

	if err := validators.ValidateID(uriParams.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
	}

	err := sc.shootService.DeleteShoot(c.Request.Context(), uriParams.ID)
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

