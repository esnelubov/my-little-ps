package pool

import (
	"errors"
	"fmt"
	"my-little-ps/common/config"
	"sync"
	"time"
)

type TaskPool struct {
	tokens   chan struct{}
	shutdown chan struct{}
	wg       *sync.WaitGroup
}

func New(config config.IConfig) *TaskPool {
	return &TaskPool{
		tokens:   make(chan struct{}, config.GetInt64("maxPoolTasks")),
		shutdown: make(chan struct{}),
		wg:       &sync.WaitGroup{},
	}
}

func (t *TaskPool) RunTask(fn func()) {
	go func() {
		select {
		case t.tokens <- struct{}{}:
			t.wg.Add(1)
			defer func() { <-t.tokens; t.wg.Done() }()
			fn()
		case <-t.shutdown:
			t.shutdown <- struct{}{}
			return
		}
	}()
}

func (t *TaskPool) GracefulShutdown(timeout time.Duration) error {
	t.shutdown <- struct{}{}

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
