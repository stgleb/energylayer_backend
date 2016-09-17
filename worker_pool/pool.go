package worker_pool

type WorkerPool interface {
	Submit(job Job)
}
