package worker_pool

// Worker
type Worker struct {
	Input chan chan Job
	Stop  chan struct{}
}

func (worker Worker) Run() {
	for {
		select {
		case jobChan := <-worker.Input:
			job := <-jobChan
			job.Do()
		}
	}
}
