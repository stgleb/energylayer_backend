package worker_pool

import (
	. "../storage"
)

const BUFFER_SIZE = 1024

// Worker
type Worker struct {
	Pool           FixedPool
	ReceiveChannel chan []Measurement
	Stop           chan struct{}
	DBStorage      Storage
}

func NewWorker(pool FixedPool, stop chan struct{}, storage Storage) Worker {
	return Worker{
		Pool:           pool,
		ReceiveChannel: make(chan Measurement),
		Stop:           stop,
		DBStorage:      storage,
	}
}

func (worker Worker) Run() {
	go func() {

		for {
			// Send workers receive channel to pool, to get new jobs.
			worker.Pool.JobQueue() <- worker.ReceiveChannel

			select {
			case measurements := <-worker.ReceiveChannel:
			// Save data to storage
				worker.DBStorage.CreateMeasurements(measurements)
			case <-worker.Stop:
				// Stop the worker
				return
			}
		}
	}()
}

func (worker Worker) Stop() {
	worker.Stop <- struct{}{}
}
