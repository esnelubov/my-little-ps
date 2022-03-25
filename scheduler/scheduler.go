package scheduler

import (
	"context"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"my-little-ps/config"
	"my-little-ps/logger"
	"reflect"
	"runtime"
	"time"
)

type Scheduler struct {
	logger *logger.Log
	cron   *cron.Cron
}

func New(logger *logger.Log, config config.IConfig) *Scheduler {
	location, err := time.LoadLocation(config.GetString("location"))
	if err != nil {
		log.Fatalf("failed to get the current location: %s", err)
	}

	return &Scheduler{
		logger: logger,
		cron:   cron.New(cron.WithLocation(location)),
	}
}

func (s *Scheduler) Start() {
	s.logger.Debug("Starting scheduler")
	s.cron.Start()
}

func (s *Scheduler) GracefulShutdown(timeout time.Duration) error {
	s.logger.Debugf("Stopping scheduler gracefully, with timeout: %v", timeout)

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
	s.logger.Debugf("Adding task: %s with schedule: %s", runtime.FuncForPC(reflect.ValueOf(cmd).Pointer()).Name(), spec)

	if _, err := s.cron.AddFunc(spec, cmd); err != nil {
		log.Fatal(err)
	}
}
