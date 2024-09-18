package main

import (
	"sync"
)

// Create a worker thread pool of 16 threads each taking input from a channel
// This is used to process the blocks concurrently
type WorkerPool struct {
	WorkerCount int
	JobQueue    chan *Job
	WaitGroup   sync.WaitGroup
	Handler     JobHandler
}

// Create a new worker pool
func NewWorkerPool(workerCount int, handler JobHandler) *WorkerPool {
	return &WorkerPool{
		WorkerCount: workerCount,
		JobQueue:    make(chan *Job, workerCount),
		Handler:     handler,
	}
}

// Start the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.WorkerCount; i++ {
		wp.WaitGroup.Add(1)
		go wp.worker(i)
	}
}

// Stop the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.JobQueue)
}

// Wait for the worker pool to finish
func (wp *WorkerPool) Wait() {
	wp.WaitGroup.Wait()
}

// Add a job to the worker pool
func (wp *WorkerPool) AddJob(job *Job) {
	wp.JobQueue <- job
}

// Worker function which processes the job
func (wp *WorkerPool) worker(workerId int) {
	defer wp.WaitGroup.Done()

	for job := range wp.JobQueue {
		wp.Handler.Process(workerId, job)
	}
}
