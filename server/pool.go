package server

type Pool struct {
	// Registered workers to the pool
	workers map[*Worker]bool

	// Inbound jobs from workers
	dispatch chan []byte

	// Register the requests coming from the workers
	register chan *Worker

	// Unregister the requests from workers
	unregister chan *Worker
}

func NewPool() *Pool {

	return &Pool{
		workers:    make(map[*Worker]bool),
		dispatch:   make(chan []byte),
		register:   make(chan *Worker),
		unregister: make(chan *Worker),
	}
}

func (pool *Pool) Run() {
	isRunning := true

	for isRunning {
		select {
		case worker := <-pool.register:
			pool.workers[worker] = true

		case worker := <-pool.unregister:
			if _, ok := pool.workers[worker]; ok {
				delete(pool.workers, worker)
				close(worker.payload)
			}

		case job := <-pool.dispatch:
			for worker := range pool.workers {
				select {
				case worker.payload <- job:
				default:
					delete(pool.workers, worker)
					close(worker.payload)
				}
			}
		}

	}
}
