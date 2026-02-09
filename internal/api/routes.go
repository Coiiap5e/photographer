package api

import (
	"github.com/Coiiap5e/photographer/internal/api/controllers"
	"github.com/gin-gonic/gin"
)

func SetupShootRoutes(router *gin.RouterGroup, shootController *controllers.ShootController) {
	shoots := router.Group("/shoots")
	{
		shoots.POST("", shootController.CreateShoot)
		shoots.GET("/:id", shootController.GetShootByID)
		shoots.GET("", shootController.GetShoots)
		shoots.DELETE("/:id", shootController.DeleteShoot)
		shoots.PUT("/:id", shootController.UpdateShoot)
		shoots.PATCH("/:id/datetime", shootController.UpdateShootDateTime)
	}
}

func SetupClientRoutes(router *gin.RouterGroup, clientController *controllers.ClientController) {
	clients := router.Group("/clients")
	{
		clients.POST("", clientController.CreateClient)
	}
}

