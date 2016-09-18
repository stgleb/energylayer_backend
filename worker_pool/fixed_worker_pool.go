package worker_pool

type FixedPool struct {
	queue           chan chan Job
	input           chan Job
	stop            chan struct{}
	stopChannels    []chan struct{}
	workerCount     int
	maxWorkersCount int
}

// NOTE: maximum quantity of goroutines is max(queueSize, maxWorkersCount)
func NewWorkerPool(queueSize, maxWorkersCount int) {
	pool := FixedPool{
		queue:           make(chan chan Job, maxWorkersCount),
		input:           make(chan Job, queueSize),
		stop:            make(chan struct{}),
		stopChannels:    make(chan struct{}),
		maxWorkersCount: maxWorkersCount,
		workerCount:     0,
	}
	// Start worker pool event-loop
	go pool.run()

	return pool
}

func (pool FixedPool) Submit(job Job) {
	pool.input <- job
}

func (pool FixedPool) JobQueue() chan chan Job {
	return pool.queue
}

func (pool FixedPool) Stop() {
	pool.stop <- struct{}{}
}

func (pool FixedPool) run() {
	for {
		select {
		case job := <-pool.input:
			// If worker limit is not exceed - spawn new worker
			if pool.workerCount < pool.maxWorkersCount {
				stopChannel := make(chan struct{})
				worker := NewWorker(pool, stopChannel)
				pool.stopChannels = append(pool.stopChannels,
					stopChannel)
				go worker.Run()
			}
			// Obtain input channel of worker
			jobChan := <-pool.queue
			// Submit job to worker
			jobChan <- job
		case <-pool.stop:
			for stopChannel := range pool.stopChannels {
				stopChannel <- struct{}{}
			}
			return
		}
	}
}
