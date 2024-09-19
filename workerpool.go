package main

import (
	"log"
	"sync"
)

// Create a worker thread pool of 16 threads each taking input from a channel
// This is used to process the blocks concurrently
type WorkerPool struct {
	WorkerCount int
	JobQueue    chan *Job
	WaitGroup   sync.WaitGroup
	Callback    func(int, *Job) error
}

// Create a new worker pool
func NewWorkerPool(workerCount int, callback func(int, *Job) error) *WorkerPool {
	return &WorkerPool{
		WorkerCount: workerCount,
		JobQueue:    make(chan *Job, workerCount),
		Callback:    callback,
	}
}

// Start the worker pool
func (wp *WorkerPool) Start() {
	if wp.Callback == nil {
		log.Println("error: no callback function")
		return
	}

	for i := 0; i < wp.WorkerCount; i++ {
		wp.WaitGroup.Add(1)
		go wp.worker(i)
	}
}

// Stop the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.JobQueue)
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
		err := wp.Callback(workerId, job)
		if err != nil {
			log.Println("Worker", workerId, "error processing job", job.BlockIndex, ":", err)
		}
	}
}
