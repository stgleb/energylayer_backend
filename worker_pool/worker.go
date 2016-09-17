package worker_pool

// Worker
type Worker struct {
	Pool           WorkerPool
	ReceiveChannel chan Job
	Stop           chan struct{}
}

func NewWorker(pool WorkerPool, stop chan struct{}) Worker {
	return Worker{
		Pool:           pool,
		ReceiveChannel: make(chan Job),
		Stop:           stop,
	}
}

func (worker Worker) Run() {
	go func() {
		for {
			// Send workers receive channel to pool, to get new jobs.
			worker.Pool.JobQueue() <- worker.ReceiveChannel

			select {
			case job := <-worker.ReceiveChannel:
				job.Do()
			case <-worker.Stop:
				// Stop the worker
				return
			}
		}
	}()
}

func (worker Worker) Stop() {
	go func() {
		worker.Stop <- struct{}{}
	}()
}
