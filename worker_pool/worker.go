package worker_pool

// Worker
type Worker struct {
	Input chan chan Job
	Stop  chan struct{}
}

func NewWorker(input chan chan struct{}, stop chan struct{}) Worker {
	return Worker{
		Input: input,
		Stop:  stop,
	}
}

func (worker Worker) Run() {
	for {
		jobChan := <-worker.Input
		job := <-jobChan
		job.Do()
	}
}
