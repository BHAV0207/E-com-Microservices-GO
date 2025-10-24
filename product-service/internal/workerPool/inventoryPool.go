package workerpool

import (
	"fmt"
	"sync"
)

type Job struct {
	ProductId string
	Action    func() error
}

type WorkerPool struct {
	Jobs    chan Job
	Workers int
	wg      sync.WaitGroup
}

func NewWorkerPool(workerCount int) *WorkerPool {
	pool := &WorkerPool{
		Jobs:    make(chan Job, 100),
		Workers: workerCount,
	}

	pool.start()
	return pool
}

func (wp *WorkerPool) start() {
	for i := 0; i < wp.Workers; i++ {
		wp.wg.Add(1)
		go func(workerId int) {
			defer wp.wg.Done()
			for job := range wp.Jobs {
				fmt.Printf("👷 Worker %d processing product: %s\n", workerId, job.ProductId)
				if err := job.Action(); err != nil {
					fmt.Printf("⚠️ Worker %d failed for %s: %v\n", workerId, job.ProductId, err)
				}
			}
		}(i + 1)
	}
}

// Submit pushes a job to the worker pool
func (wp *WorkerPool) Submit(job Job) {
	wp.Jobs <- job
}

// Shutdown gracefully closes the pool
func (wp *WorkerPool) Shutdown() {
	close(wp.Jobs)
	wp.wg.Wait()
}
