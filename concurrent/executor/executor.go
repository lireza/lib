// Package executor defines the thread pool executor abstraction.
// It also contains some predefined executors.
package executor

import (
	"container/ring"
	"errors"
	"runtime"
	"sync"

	"github.com/lireza/lib/concurrent"
)

// Executor is the thread pool abstraction.
type Executor interface {
	// Execute executes the runner instance passed to method.
	// Whether or not the runner will be called on new thread depends on implementation.
	Execute(runner concurrent.Runner)

	// Shutdown shutdowns the executor and waits on threads to complete their tasks
	// passed to threads before the shutdown signal.
	Shutdown()
}

// RoundRobinExecutor is an executor implementation that contains some threads,
// and passes tasks to threads in a round robin fashion.
type RoundRobinExecutor struct {
	mutex    *sync.Mutex
	ids      *ring.Ring
	channels map[int]chan concurrent.Runner
	shutdown map[int]chan struct{}
	wg       *sync.WaitGroup
}

// NewRoundRobinExecutor creates a new executor based on round robin distribution concept.
// The number of threads in executor is defined through nThreads.
// Each thread will have a queue for runners and the size of queues is defined through threadQueueSize.
// In case of errors during executor creation the error will be return.
func NewRoundRobinExecutor(nThreads, threadQueueSize int) (*RoundRobinExecutor, error) {
	if nThreads < 1 || threadQueueSize < 1 {
		return nil, errors.New("executor: invalid argument")
	}

	ids := ring.New(nThreads)
	for i := 1; i <= nThreads; i++ {
		ids.Value = i
		ids = ids.Next()
	}

	channels := make(map[int]chan concurrent.Runner, nThreads)
	shutdown := make(map[int]chan struct{}, nThreads)
	wg := &sync.WaitGroup{}
	wg.Add(nThreads)

	for i := 1; i <= nThreads; i++ {
		channels[i] = make(chan concurrent.Runner, threadQueueSize)
		shutdown[i] = make(chan struct{})

		go func(id int, runners <-chan concurrent.Runner, shutdown <-chan struct{}, wg *sync.WaitGroup) {
			runtime.LockOSThread()

			for {
				select {
				case runner := <-runners:
					runner.Run()
				case <-shutdown:
					runtime.UnlockOSThread()
					wg.Done()
					return
				}
			}
		}(i, channels[i], shutdown[i], wg)
	}

	return &RoundRobinExecutor{mutex: &sync.Mutex{}, ids: ids, channels: channels, shutdown: shutdown, wg: wg}, nil
}

func (e *RoundRobinExecutor) Execute(runner concurrent.Runner) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	id := e.ids.Value.(int)
	e.channels[id] <- runner
	e.ids = e.ids.Next()
}

func (e *RoundRobinExecutor) Shutdown() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, c := range e.shutdown {
		c <- struct{}{}
	}

	e.wg.Wait()
}
