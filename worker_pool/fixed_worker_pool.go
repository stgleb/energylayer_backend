package worker_pool

type WorkerPool struct {
	queue           chan chan Job
	input           chan Job
	Stop            chan struct{}
	maxWorkersCount int
}

// NOTE: maximum quantity of gorotines is max(queueSize, maxWorkersCount)
func NewWorkerPool(queueSize, maxWorkersCount int) {
	pool := WorkerPool{
		queue:           make(chan chan Job, maxWorkersCount),
		input:           make(chan Job, queueSize),
		Stop:            make(chan struct{}),
		maxWorkersCount: maxWorkersCount,
	}
	pool.run()

	return pool
}

func (pool WorkerPool) Submit(job Job) {
	pool.input <- job
}

func (pool WorkerPool) JobQueue() chan chan Job {
	return pool.queue
}

func (pool WorkerPool) run() {
	for {
		select {
		case job := <-pool.input:
			// TODO: remove this goroutine spawn in flavor of
			// blocking call
			go func(job Job) {
				// Obtain input channel of worker
				jobChan := <-pool.queue
				// Submit job to worker
				jobChan <- job
			}(job)
		case <-pool.Stop:
			return
		}
	}
}
