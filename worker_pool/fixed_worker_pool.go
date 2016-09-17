package worker_pool

type WorkerPool struct {
	Queue           chan chan Job
	Input           chan Job
	Stop            chan struct{}
	maxWorkersCount int
}

// NOTE: maximum quantity of gorotines is max(queueSize, maxWorkersCount)
func NewWorkerPool(queueSize, maxWorkersCount int) {
	pool := WorkerPool{
		Queue:           make(chan chan Job, maxWorkersCount),
		Input:           make(chan Job, queueSize),
		Stop:            make(chan struct{}),
		maxWorkersCount: maxWorkersCount,
	}
	pool.run()

	return pool
}

func (pool WorkerPool) Submit(job Job) {
	pool.Queue <- job
}

func (pool WorkerPool) run() {
	for {
		select {
		case job := <-pool.Input:
			// TODO: remove this goroutine spawn in flavor of
			// blocking call
			go func(job Job) {
				// Obtain input channel of worker
				jobChan := <-pool.Queue
				// Submit hob to worker
				jobChan <- jobChan
			}(job)
		case <-pool.Stop:
			return
		}
	}
}
