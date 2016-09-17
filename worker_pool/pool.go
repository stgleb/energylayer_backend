package worker_pool

type WorkerPool interface {
	Submit(job Job)
	JobQueue() chan chan Job
}
