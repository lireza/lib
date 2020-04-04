// Package concurrent provides some utility abstractions and functions that are used in concurrent programming.
package concurrent

// Runner is an abstraction for an execution that can be start with calling run method.
// Runners can be passed to goroutines, so should be concurrent safe.
type Runner interface {
	// Run will be called on runner instances to start the task defined.
	Run()
}

// Task is a function and has a channel of interface to send the response of function execution.
// It also accepts an argument to be passed to the function.
type Task struct {
	do  func(interface{}, chan<- interface{})
	arg interface{}
	r   chan<- interface{}
}

func (t *Task) Run() {
	t.do(t.arg, t.r)
}

// NewTask creates a new task and also returns the response channel to wait on.
func NewTask(do func(interface{}, chan<- interface{}), arg interface{}) (*Task, <-chan interface{}) {
	// To protect the task invoker's goroutine from blocking.
	r := make(chan interface{}, 2)
	return &Task{do: do, arg: arg, r: r}, r
}
