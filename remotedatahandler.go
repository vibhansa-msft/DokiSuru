package main

import (
	"dokisuru/storage"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

var _ JobHandler = &RemoteDataHandler{}

type RemoteDataHandler struct {
	BaseHandler
	Container *container.Client
}

func NewRemoteDataHandler(workerCount int, next JobHandler) *RemoteDataHandler {
	rdh := &RemoteDataHandler{
		BaseHandler: BaseHandler{
			Next: next},
	}
	rdh.Worker = NewWorkerPool(workerCount, rdh.Process)
	return rdh
}

// Process the block
func (rdh *RemoteDataHandler) Process(workerId int, bj *Job) error {
	// Open the file
	file, err := os.Open(bj.Path)
	if err != nil {
		log.Println("Worker", workerId, "error opening file", bj.Path, ":", err)
		return fmt.Errorf("error opening file %s: %v", bj.Path, err)
	}

	defer file.Close()

	// Create the block index
	bj.BlockIndex = uint16(bj.Offset / int64(config.BlockSize))

	// Read the data from given offset
	bj.Data = make([]byte, config.BlockSize)
	n, err := file.ReadAt(bj.Data, bj.Offset)
	if err != nil {
		if err != io.EOF {
			log.Println("Worker", workerId, "error reading file", bj.Path, ":", err)
			return fmt.Errorf("error reading file %s: %v", bj.Path, err)
		}
	}
	if n <= 0 {
		log.Println("Worker", workerId, "nothing to process from file", bj.Path, ":", err)
		return fmt.Errorf("nothing to process from file %s: %v", bj.Path, err)
	}

	// Compuate md5sum of the data
	bj.Md5Sum = ComputeMd5Sum(bj.Data[:n])

	// Convert this slice to a base64 encoded string
	bj.BlockId = GetBlockID(bj.BlockIndex, bj.Md5Sum)

	log.Println("Worker", workerId, "processed block", bj.BlockIndex, "with blockId", bj.BlockId)
	log.Printf("%x\n", bj.Md5Sum)

	return nil
}

func (rdh *RemoteDataHandler) Start() error {
	var err error
	clients, err := storage.NewClients()
	if err != nil {
		log.Println("Error creating clients:", err)
		return err
	}

	rdh.Container = clients.GetContainerClient()

	return nil
}
