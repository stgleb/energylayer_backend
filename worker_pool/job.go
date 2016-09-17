package worker_pool

// Job interface representing unit of work,
// do function contains task to do, Cancel return
// write only channel to cancel the task execution.
type Job interface {
	Do()
	Cancel() chan<- struct{}
}
