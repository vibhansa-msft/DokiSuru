package main

import (
	"flag"
	"log"
	_ "net/http/pprof"
	"os"
)

func main() {
	flag.IntVar(&config.WorkerCount, "worker", 16, "Number of workers to use")
	flag.Uint64Var(&config.BlockSize, "blocksize", 16*1024*1024, "Block size to use")
	flag.StringVar(&config.Path, "path", "testdata.1g", "Path to the file to process")

	file, err := os.OpenFile("dokisuru.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("DōkiSuru: Starting up")

	// Parse the user config
	flag.Parse()

	// Create a local data handler
	localDataHandler := NewLocalDataHandler()

	// Start the worker pool
	localWorkers := NewWorkerPool(config.WorkerCount, localDataHandler)
	localWorkers.Start()

	localDataHandler.Start(localWorkers.AddJob)

	// Stop the worker pool
	localWorkers.Stop()

	// Wait for the localWorkers to finish
	localWorkers.Wait()

	log.Printf("DōkiSuru: Finishing")
}
