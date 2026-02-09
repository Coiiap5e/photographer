package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Coiiap5e/photographer/internal/api"
	"github.com/Coiiap5e/photographer/internal/api/middleware"
	"github.com/Coiiap5e/photographer/internal/app/di"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	container, err := di.NewContainer(ctx)
	if err != nil {
		log.Fatal("Failed to create DI container:", err)
	}

	defer container.Close()

	container.Scheduler.Start()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(container.Logger))
	router.Use(middleware.ErrorHandler())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	apiGroup := router.Group("/api")
	api.SetupShootRoutes(apiGroup, container.Controllers.Shoot)
	api.SetupClientRoutes(apiGroup, container.Controllers.Client)

	addr := fmt.Sprintf("%s:%d",
		container.Config.Server.Host,
		container.Config.Server.Port,
	)

	container.Logger.Info("Starting server", "address", addr)

	if err = router.Run(addr); err != nil {
		container.Logger.Error("Failed to start server", "error", err)
		log.Fatal("Failed to start server:", err)
	}
}
