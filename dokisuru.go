package main

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
)

func main() {
	flag.IntVar(&config.WorkerCount, "worker", 16, "Number of workers to use")
	flag.Uint64Var(&config.BlockSize, "blocksize", 16*1024*1024, "Block size to use")
	flag.StringVar(&config.Path, "path", "./README.md", "Path to the file to process")

	file, err := os.OpenFile("dokisuru.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("DōkiSuru: Starting up")
	log.Printf("DōkiSuru: Starting up")

	// Parse the user config
	flag.Parse()

	// Start the worker pool
	workerPool := NewWorkerPool(config.WorkerCount)
	workerPool.Start()

	// Create a local data handler
	localDataHandler := NewLocalDataHandler(workerPool)

	localDataHandler.Start()

	// Stop the worker pool
	workerPool.Stop()

	// Wait for the workerpool to finish
	workerPool.Wait()

	fmt.Println("DōkiSuru: Finishing")
	log.Printf("DōkiSuru: Finishing")
}
