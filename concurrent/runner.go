// Package concurrent provides some utility abstractions and functions that are used in concurrent programming.
package concurrent

// Runner is an abstraction for an execution that can be start with calling run method.
// Runners can be passed to goroutines and so should be concurrent safe.
type Runner interface {
	// Run will be called on runnable instances to start the work defined.
	Run()
}

// Callable has a function to do its work and a channel of error to send the error occurred during work execution.
// If the receiver received the nil error it means the callable work was successful.
type Callable struct {
	do func(e chan<- error)
	e  chan<- error
}

func (c *Callable) Run() {
	c.do(c.e)
}

// NewCallable creates a new callable and also returns the error channel to wait on.
// If the receiver received the nil error in channel it means the callable work was successful.
func NewCallable(do func(e chan<- error)) (*Callable, <-chan error) {
	// To protect the callable's thread from blocking.
	e := make(chan error, 2)
	return &Callable{do: do, e: e}, e
}
