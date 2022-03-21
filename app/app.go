package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

type App struct {
	*fiber.App
}

func New(config ...fiber.Config) *App {
	return &App{fiber.New(config...)}
}

func (this *App) GracefulShutdown(timeout time.Duration) error {
	shutDown := make(chan struct{}, 1)
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), timeout)
	defer timeoutCancel()

	go func() {
		err := this.Shutdown()
		if err != nil {
			log.Fatalf("couldn't shutdown the app: %s", err)
		}

		shutDown <- struct{}{}
	}()

	select {
	case <-shutDown:
		return nil
	case <-timeoutCtx.Done():
		return errors.New(fmt.Sprintf("app didn't shutdown within the given duration: %s", timeout))
	}
}
