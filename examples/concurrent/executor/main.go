package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/lireza/lib/concurrent"
	"github.com/lireza/lib/concurrent/executor"
)

func main() {
	// Creating a new round robin executor. It distributes tasks in round robin fashion to threads.
	ex, _ := executor.NewRoundRobinExecutor(32, 500)
	wg := &sync.WaitGroup{}
	wg.Add(5_000)

	start := time.Now()
	// 5000 callable to executor. Each have 5 milliseconds of blocking code.
	for i := 1; i <= 5_000; i++ {
		c, _ := concurrent.NewCallable(func(e chan<- error) {
			time.Sleep(5 * time.Millisecond)
			wg.Done()
		})

		ex.Execute(c)
	}

	wg.Wait()
	ex.Shutdown()

	fmt.Println(time.Now().Sub(start))
}
