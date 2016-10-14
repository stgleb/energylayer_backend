package worker_pool

import (
	. "../storage"
)

type Pool interface {
	Submit(Measurement)
	Stop()
	JobQueue() chan chan Measurement
}
