package main

import (
	"log"
	"sync"
)

// Create a worker thread pool of 16 threads each taking input from a channel
// This is used to process the blocks concurrently
type WorkerPool struct {
	WorkerCount int
	JobQueue    chan Job
	WaitGroup   sync.WaitGroup
}

// Create a new worker pool
func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		WorkerCount: workerCount,
		JobQueue:    make(chan Job, workerCount),
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
func (wp *WorkerPool) AddJob(job Job) {
	wp.JobQueue <- job
}

// Worker function which processes the job
func (wp *WorkerPool) worker(workerId int) {
	log.Println("Worker", workerId, "started")
	defer wp.WaitGroup.Done()

	for job := range wp.JobQueue {
		job.Process(workerId)

		// Push the job to next queue
	}

	log.Println("Worker", workerId, "done")
}
