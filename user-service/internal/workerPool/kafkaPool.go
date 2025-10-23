package workerpool

import (
	"log"

	"github.com/BHAV0207/user-service/internal/events"
)

// why are we using pool and intead of goroutines , because we cant control the number of go routines , if 1 mill request comes then there will be 1 mill go eutines , whhc will cause cpu overload , in pool we can decide at once how many goroutines will be there working

type Job struct {
	KafkaEvent map[string]any
}

//  suppose we have more things other that kafka, we have two opitons either cereate a seperate pool or in the smae pool where we have made the Job stuct there we can pass the function named task which will be genereic function and all the events can define that function and pass that functin

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
	for i := 0; i < wp.Workers; i++ {
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
