package executor

import (
	"container/ring"
	"errors"
	"runtime"
	"sync"

	"github.com/lireza/lib/concurrent"
)

type Executor interface {
	Execute(task *concurrent.Task)
	Shutdown()
}

type RoundRobinExecutor struct {
	mutex    *sync.Mutex
	ids      *ring.Ring
	channels map[int]chan *concurrent.Task
	shutdown map[int]chan struct{}
	wg       *sync.WaitGroup
}

func NewRoundRobinExecutor(nThreads, threadQueueSize int) (*RoundRobinExecutor, error) {
	if nThreads < 1 || threadQueueSize < 1 {
		return nil, errors.New("executor: invalid argument")
	}

	ids := ring.New(nThreads)
	for i := 1; i <= ids.Len(); i++ {
		ids.Value = i
		ids = ids.Next()
	}

	channels := make(map[int]chan *concurrent.Task, nThreads)
	shutdown := make(map[int]chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(nThreads)

	for i := 1; i <= nThreads; i++ {
		channels[i] = make(chan *concurrent.Task, threadQueueSize)
		shutdown[i] = make(chan struct{})

		go func(id int, tasks <-chan *concurrent.Task, shutdown <-chan struct{}, wg *sync.WaitGroup) {
			runtime.LockOSThread()

			for {
				select {
				case t := <-tasks:
					t.Run()
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

func (e *RoundRobinExecutor) Execute(task *concurrent.Task) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	id := e.ids.Value.(int)
	e.channels[id] <- task
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
