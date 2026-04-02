package controllers

import (
	"net/http"

	"github.com/Coiiap5e/photographer/internal/api/controllers/dto"
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
	var req dto.CreateShootRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := validators.ValidateCreateShoot(&req, sc.clock); err != nil {
		_ = c.Error(err)
		return
	}

	clients := dto.ToShootClientDomain(req.Clients)

	newShoot := &model.Shoot{
		ShootDate:     req.ShootDate,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		ShootPrice:    req.ShootPrice,
		ShootLocation: req.ShootLocation,
		ShootType:     req.ShootType,
		Notes:         req.Notes,
	}

	ctx := c.Request.Context()
	createdShoot, err := sc.shootService.CreateShoot(ctx, newShoot, clients)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToShootResponse(createdShoot))
}

func (sc *ShootController) UpdateShoot(c *gin.Context) {
	var uriParams struct {
		ID int `uri:"id" binding:"min=1"`
	}
	if err := c.ShouldBindUri(&uriParams); err != nil {
		_ = c.Error(err)
		return
	}

	var req dto.CreateShootRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := validators.ValidateCreateShoot(&req, sc.clock); err != nil {
		_ = c.Error(err)
		return
	}
	if err := validators.ValidateID(uriParams.ID); err != nil {
		_ = c.Error(err)
		return
	}

	clients := dto.ToShootClientDomain(req.Clients)

	newShoot := &model.Shoot{
		ShootDate:     req.ShootDate,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		ShootPrice:    req.ShootPrice,
		ShootLocation: req.ShootLocation,
		ShootType:     req.ShootType,
		Notes:         req.Notes,
	}

	ctx := c.Request.Context()
	updatedShoot, err := sc.shootService.UpdateShoot(ctx, uriParams.ID, newShoot, clients)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToShootResponse(updatedShoot))
}

func (sc *ShootController) UpdateShootDateTime(c *gin.Context) {
	var uriParams struct {
		ID int `uri:"id" binding:"min=1"`
	}
	if err := c.ShouldBindUri(&uriParams); err != nil {
		_ = c.Error(err)
		return
	}

	var req dto.UpdateShootDateTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := validators.ValidateDate(&req, sc.clock); err != nil {
		_ = c.Error(err)
		return
	}
	if err := validators.ValidateID(uriParams.ID); err != nil {
		_ = c.Error(err)
		return
	}

	newUpdatedDateTime := &model.ShootDateTimePatch{
		ShootDate: req.ShootDate,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	ctx := c.Request.Context()
	updatedShoot, err := sc.shootService.UpdateShootDateTime(ctx, uriParams.ID, newUpdatedDateTime)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToShootResponse(updatedShoot))
}

func (sc *ShootController) GetShootByID(c *gin.Context) {
	var uriParams dto.GetShootByIDRequest
	if err := c.ShouldBindUri(&uriParams); err != nil {
		_ = c.Error(err)
		return
	}

	if err := validators.ValidateID(uriParams.ID); err != nil {
		_ = c.Error(err)
		return
	}

	ctx := c.Request.Context()
	shoot, err := sc.shootService.GetShootByID(ctx, uriParams.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToShootResponse(shoot))
}

func (sc *ShootController) GetShoots(c *gin.Context) {
	ctx := c.Request.Context()
	shoots, err := sc.shootService.GetShoots(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	responseShoots := make([]*dto.ShootResponse, 0)
	for _, shoot := range shoots {
		responseShoots = append(responseShoots, dto.ToShootResponse(&shoot))
	}

	c.JSON(http.StatusOK, responseShoots)
}

func (sc *ShootController) DeleteShoot(c *gin.Context) {
	var uriParams struct {
		ID int `uri:"id" binding:"min=1"`
	}
	if err := c.ShouldBindUri(&uriParams); err != nil {
		_ = c.Error(err)
		return
	}

	if err := validators.ValidateID(uriParams.ID); err != nil {
		_ = c.Error(err)
		return
	}

	ctx := c.Request.Context()
	err := sc.shootService.DeleteShoot(ctx, uriParams.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shoot deleted successfully",
		"deleted": true,
	})
}
