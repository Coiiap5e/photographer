package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"

	"github.com/Coiiap5e/photographer/internal/api"
	"github.com/Coiiap5e/photographer/internal/api/middleware"
	"github.com/Coiiap5e/photographer/internal/app/di"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	container, err := di.NewContainer(ctx)
	if err != nil {
		log.Fatal("failed to create DI container:", err)
	}
	defer container.Close()

	g, gCtx := errgroup.WithContext(ctx)

	container.Scheduler.Start()

	g.Go(func() error {
		<-gCtx.Done()
		container.Scheduler.Stop()
		container.Logger.Info("scheduler stopped")
		return nil
	})

	g.Go(func() error {
		return container.WorkerPool.Run(gCtx)
	})

	server := setupServer(container)
	g.Go(func() error {
		container.Logger.Info("starting server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			container.Logger.Error("failed to start server", "error", err)
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		container.Logger.Info("shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			container.Logger.Error("server shutdown failed", "error", err)
			return err
		}

		container.Logger.Info("server stopped gracefully")
		return nil
	})

	container.Logger.Info("application is running...")
	if err := g.Wait(); err != nil {
		container.Logger.Error("application finished with an error", "error", err)
	} else {
		container.Logger.Info("application finished gracefully")
	}
}

func setupServer(container *di.Container) *http.Server {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(container.Logger))
	router.Use(middleware.ErrorHandler())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	apiGroup := router.Group("/api")
	api.SetupShootRoutes(apiGroup, container.Controllers.Shoot)
	api.SetupClientRoutes(apiGroup, container.Controllers.Client)

	addr := fmt.Sprintf("%s:%d", container.Config.Server.Host, container.Config.Server.Port)
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
