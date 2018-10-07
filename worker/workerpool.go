package worker

import (
	"sync"
)

type WorkerPool struct {
	name       string
	maxworkers int
	wg         sync.WaitGroup
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.maxworkers; i++ {
		wp.wg.Add(1)
		go wp.run()
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func (wp *WorkerPool) run() {
	wp.wg.Done()
}
