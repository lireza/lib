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
	// 5000 tasks to executor. Each have 5 milliseconds of blocking code.
	for i := 1; i <= 5_000; i++ {
		t := concurrent.NewTask(func(i interface{}, c chan<- interface{}) {
			time.Sleep(5 * time.Millisecond)
			i.(*sync.WaitGroup).Done()
		}, wg, nil)

		ex.Execute(t)
	}

	wg.Wait()
	ex.Shutdown()
	// 900.447083ms on my laptop.
	fmt.Println(time.Now().Sub(start))
}
