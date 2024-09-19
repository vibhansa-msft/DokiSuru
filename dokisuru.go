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

	// Create a remote data handler
	RemoteDataHandler := NewRemoteDataHandler(config.WorkerCount, nil)

	// Create a local data handler
	localDataHandler := NewLocalDataHandler(config.WorkerCount, RemoteDataHandler)

	// Start the worker pool
	err = RemoteDataHandler.Start()
	if err != nil {
		fmt.Println("Error starting remote data handler: %v", err)
		return
	}

	err = localDataHandler.Start()
	if err != nil {
		fmt.Println("Error starting local data handler: %v", err)
		return
	}

	// Start the worker pool
	localDataHandler.Stop()
	RemoteDataHandler.Stop()

	log.Printf("DōkiSuru: Finishing")
}
