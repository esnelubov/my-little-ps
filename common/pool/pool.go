package pool

import (
	"errors"
	"fmt"
	"my-little-ps/common/config"
	"my-little-ps/common/logger"
	"sync"
	"time"
)

type TaskPool struct {
	logger   *logger.Log
	tokens   chan struct{}
	shutdown chan struct{}
	wg       *sync.WaitGroup
	Error    error
}

func New(logger *logger.Log, config config.IConfig) *TaskPool {
	return &TaskPool{
		logger:   logger,
		tokens:   make(chan struct{}, config.GetInt64("maxPoolTasks")),
		shutdown: make(chan struct{}, config.GetInt64("maxPoolTasks")),
		wg:       &sync.WaitGroup{},
	}
}

func (t *TaskPool) RunTask(fn func() error) {
	t.wg.Add(1)

	go func() {
		select {
		case t.tokens <- struct{}{}:
			defer func() { <-t.tokens; t.wg.Done() }()
			err := fn()
			if err != nil {
				t.Error = err
				t.Shutdown()
			}
		case <-t.shutdown:
			t.shutdown <- struct{}{}
			return
		}
	}()
}

func (t *TaskPool) WaitTasks() *TaskPool {
	t.wg.Wait()
	return t
}

func (t *TaskPool) Shutdown() {
	t.shutdown <- struct{}{}
}

func (t *TaskPool) GracefulShutdown(timeout time.Duration) error {
	t.logger.Debugf("Stopping pool gracefully, with timeout: %v", timeout)

	t.Shutdown()

	if !waitTimeout(t.wg, timeout) {
		return errors.New(fmt.Sprintf("pool tasks didn't end within the given duration: %s", timeout))
	}

	return nil
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
