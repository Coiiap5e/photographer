package routes

import (
	"github.com/Coiiap5e/photographer/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupShootRoutes(router *gin.RouterGroup, shootHandler *handlers.ShootHandler) {
	shoots := router.Group("/shoots")
	{
		shoots.POST("", shootHandler.CreateShoot)
		shoots.GET("/:id", shootHandler.GetShootByID)
		shoots.DELETE("/:id", shootHandler.DeleteShoot)
		shoots.PUT("/:id", shootHandler.UpdateShoot)
		shoots.PATCH("/:id/datetime", shootHandler.UpdateShootDateTime)
	}
}

func SetupClientRoutes(router *gin.RouterGroup, clientHandler *handlers.ClientHandler) {
	clients := router.Group("/clients")
	{
		clients.POST("", clientHandler.CreateClient)
	}
}
