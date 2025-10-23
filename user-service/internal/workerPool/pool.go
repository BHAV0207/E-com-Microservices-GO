package workerpool

import (
	"log"

	"github.com/BHAV0207/user-service/internal/events"
)

type Job struct {
	KafkaEvent map[string]any
}

type WorkerPool struct {
	Jobs     chan Job
	Workers  int
	Producer *events.Producer
}

func NewWorkerPool(workerCount int, producer *events.Producer) *WorkerPool {
	pool := &WorkerPool{
		Jobs:     make(chan Job, 1000),
		Workers:  workerCount,
		Producer: producer,
	}

	pool.start()
	return pool
}

func (wp *WorkerPool) start() {
	for i := range wp.Workers {
		go func(id int) {
			for job := range wp.Jobs {
				if err := wp.Producer.Publish(job.KafkaEvent); err != nil {
					log.Printf("⚠️ Worker %d: failed to publish event: %v", id, err)
				} else {
					log.Printf("✅ Worker %d: event published successfully", id)
				}
			}
		}(i)
	}
}

func (wp *WorkerPool) Submit(event map[string]interface{}) {
	wp.Jobs <- Job{KafkaEvent: event}
}
