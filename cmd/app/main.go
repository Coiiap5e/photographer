package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	app, cleanup, err := InitializeApp()
	if err != nil {
		fmt.Println("Error initializing app", err)
		os.Exit(1)
	}

	defer cleanup()

	logger := app.GetLogger()
	logger.Info("app starting")

	signalChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan

		fmt.Println("Got signal: ", sig)
		logger.Info("got signal", "signal", sig.String())

		done <- true
	}()

	go func() {
		app.RunMenu(ctx)
		done <- true
	}()

	<-done

	logger.Info("app end")
}
