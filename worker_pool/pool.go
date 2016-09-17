package worker_pool

type Pool interface {
	Submit(job Job)
	Stop()
	JobQueue() chan chan Job
}
