package worker_pool

import (
	. "../storage"
)

const BUFFER_SIZE = 1024

// Worker
type Worker struct {
	Pool           FixedPool
	ReceiveChannel chan Measurement
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
		buffer := make([]Measurement, 0, BUFFER_SIZE)

		for {
			// Send workers receive channel to pool, to get new jobs.
			worker.Pool.JobQueue() <- worker.ReceiveChannel

			select {
			case measurement := <-worker.ReceiveChannel:
				buffer := append(buffer, measurement)
				worker.DBStorage.CreateMeasurements(buffer)
				buffer = make([]Measurement, 0, BUFFER_SIZE)
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
