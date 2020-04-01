package concurrent

type Task struct {
	function func(interface{}, chan<- interface{})
	arg      interface{}
	response chan interface{}
}

func (t *Task) Run() {
	t.function(t.arg, t.response)
}

func NewTask(function func(interface{}, chan<- interface{}), arg interface{}, response chan interface{}) *Task {
	return &Task{function: function, arg: arg, response: response}
}
