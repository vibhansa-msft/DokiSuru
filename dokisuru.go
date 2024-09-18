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
	localDataHandler := NewLocalDataHandler(config.WorkerCount, nil)

	// Start the worker pool
	localDataHandler.Worker.Start()

	localDataHandler.Start()

	// Stop the worker pool
	localDataHandler.Worker.Stop()

	// Wait for the localWorkers to finish
	localDataHandler.Worker.Wait()

	log.Printf("DōkiSuru: Finishing")
}
