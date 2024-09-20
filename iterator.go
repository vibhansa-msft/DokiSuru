package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type Iterator struct {
	BaseHandler
	Path string
}

func NewIterator(workerCount int) *Iterator {
	it := &Iterator{
		BaseHandler: BaseHandler{
			Next: nil},
	}
	it.Worker = NewWorkerPool(workerCount, it.Process)
	it.Worker.Start()
	return it
}

func (it *Iterator) Start(path string) error {
	it.Path = path
	info, err := os.Lstat(it.Path)
	if err != nil {
		log.Println("Error getting file info for", it.Path, ":", err)
		return err
	}

	if info.IsDir() {
		log.Println("Input Path is a directory")
		err := filepath.Walk(it.Path, func(path string, info fs.FileInfo, err error) error {
			if err == nil && info != nil {
				job := Job{
					Path: path,
				}
				if !info.IsDir() {
					log.Println("Process :", path)
					it.Worker.AddJob(&job)
				}
			}
			return nil
		})
		if err != nil {
			log.Println("Error walking path", it.Path, ":", err)
			return err
		}
	} else {
		job := Job{
			Path: path,
		}

		it.Worker.AddJob(&job)
	}

	return nil
}

func (it *Iterator) Stop() error {
	it.Worker.Stop()
	return nil
}

func (it *Iterator) Process(workerId int, bj *Job) error {
	// Create a remote data handler
	RemoteDataHandler := NewRemoteDataHandler(config.WorkerCount, nil)

	// Create a local data handler
	localDataHandler := NewLocalDataHandler(config.WorkerCount, RemoteDataHandler)

	// Start the worker pool
	err := RemoteDataHandler.Start(bj.Path)
	if err != nil {
		fmt.Println("Error starting remote data handler: %v", err)
		return err
	}

	err = localDataHandler.Start(bj.Path)
	if err != nil {
		fmt.Println("Error starting local data handler: %v", err)
		return err
	}

	// Start the worker pool
	localDataHandler.Stop()
	RemoteDataHandler.Stop()

	return nil
}
