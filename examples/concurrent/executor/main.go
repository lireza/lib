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
	// Sending 5000 tasks to executor. Each have 5 milliseconds of blocking code.
	for i := 1; i <= 5_000; i++ {
		t, _ := concurrent.NewTask(func(arg interface{}, r chan<- interface{}) {
			time.Sleep(5 * time.Millisecond)
			arg.(*sync.WaitGroup).Done()
		}, wg)

		ex.Execute(t)
	}

	wg.Wait()
	ex.Shutdown()
	ex.AwaitTermination()

	fmt.Println(time.Now().Sub(start))
}
