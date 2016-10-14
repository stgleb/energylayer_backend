package worker_pool

import (
	. "../storage"
)

type FixedPool struct {
	queue           chan chan Measurement
	input           chan Measurement
	stop            chan struct{}
	stopChannels    []chan struct{}
	workerCount     int
	maxWorkersCount int
}

// NOTE: maximum quantity of goroutines is max(queueSize, maxWorkersCount)
func NewWorkerPool(queueSize, maxWorkersCount int) {
	pool := FixedPool{
		queue:           make(chan chan Measurement, maxWorkersCount),
		input:           make(chan Measurement, queueSize),
		stop:            make(chan struct{}),
		stopChannels:    make(chan struct{}),
		maxWorkersCount: maxWorkersCount,
		workerCount:     0,
	}
	// Start worker pool event-loop
	go pool.run()

	return pool
}

func (pool FixedPool) Submit(measurement Measurement) {
	pool.input <- measurement
}

func (pool FixedPool) JobQueue() chan chan Measurement {
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
				// TODO: Add storage to worker
				worker := NewWorker(pool, stopChannel, nil)
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
