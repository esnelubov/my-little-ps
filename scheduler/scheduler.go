package scheduler

import (
	"context"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"my-little-ps/config"
	"time"
)

type Scheduler struct {
	cron *cron.Cron
}

func New(config config.IConfig) *Scheduler {
	location, err := time.LoadLocation(config.GetString("location"))
	if err != nil {
		log.Fatalf("failed to get the current location: %s", err)
	}

	return &Scheduler{
		cron: cron.New(cron.WithLocation(location)),
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) GracefulShutdown(timeout time.Duration) error {
	cronCtx := s.cron.Stop()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), timeout)
	defer timeoutCancel()

	select {
	case <-cronCtx.Done():
		return nil
	case <-timeoutCtx.Done():
		return errors.New(fmt.Sprintf("cron tasks didn't end within the given duration: %s", timeout))
	}
}

func (s *Scheduler) AddTask(spec string, cmd func()) {
	if _, err := s.cron.AddFunc(spec, cmd); err != nil {
		log.Fatal(err)
	}
}
