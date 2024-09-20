package main

import (
	"dokisuru/utils"
	"fmt"
	"io"
	"log"
	"os"
)

var _ JobHandler = &LocalDataHandler{}

type LocalDataHandler struct {
	BaseHandler
}

func NewLocalDataHandler(workerCount int, next JobHandler) *LocalDataHandler {
	ldh := &LocalDataHandler{
		BaseHandler: BaseHandler{
			Next: next},
	}
	ldh.Worker = NewWorkerPool(workerCount, ldh.Process)
	ldh.Worker.Start()
	return ldh
}

func (ldh *LocalDataHandler) Start() error {
	info, err := os.Lstat(config.Path)
	if err != nil {
		log.Println("Error getting file info for", config.Path, ":", err)
		return err
	}

	if info.IsDir() {
		log.Println("Error: Path is a directory")
		return fmt.Errorf("path is a directory")
	}

	// Get the file size
	fileSize := info.Size()

	// Calculate the number of blocks
	blockCount := uint16(fileSize / int64(config.BlockSize))
	if fileSize%int64(config.BlockSize) != 0 {
		blockCount++
	}

	if blockCount > 50000 {
		log.Println("Error: File too big")
		return fmt.Errorf("file too big")
	}

	if blockCount == 0 {
		ldh.Next.Enqueue(&Job{
			BlockIndex: 0,
			NoOfBlocks: 0,
		})
	}

	// Create a job for each block
	for i := uint16(0); i < blockCount; i++ {
		job := Job{
			BlockIndex: i,
			NoOfBlocks: blockCount,
		}

		ldh.Worker.AddJob(&job)
	}
	return nil
}

func (ldh *LocalDataHandler) Stop() error {
	ldh.Worker.Stop()
	return nil
}

// Process the block
func (ldh *LocalDataHandler) Process(workerId int, bj *Job) error {
	// Open the file
	file, err := os.Open(config.Path)
	if err != nil {
		log.Println("Worker", workerId, "error opening file", config.Path, ":", err)
		return fmt.Errorf("error opening file %s: %v", config.Path, err)
	}

	defer file.Close()

	// Read the data from given offset
	bj.Data = make([]byte, config.BlockSize)
	n, err := file.ReadAt(bj.Data, int64(bj.BlockIndex)*int64(config.BlockSize))
	if err != nil {
		if err != io.EOF {
			log.Println("Worker", workerId, "error reading file", config.Path, ":", err)
			return fmt.Errorf("error reading file %s: %v", config.Path, err)
		}
	}

	if int(bj.BlockIndex) < (int(bj.NoOfBlocks)-1) && n < int(config.BlockSize) {
		log.Println("Worker", workerId, "nothing to process from file", config.Path, ":", err)
		return fmt.Errorf("nothing to process from file %s: %v", config.Path, err)
	}

	// Compuate md5sum of the data
	if int(bj.BlockIndex) == (int(bj.NoOfBlocks) - 1) {
		bj.Data = bj.Data[:n]
	}
	bj.Md5Sum = utils.ComputeMd5Sum(bj.Data)

	// Convert this slice to a base64 encoded string
	bj.BlockId = utils.GetBlockID(bj.BlockIndex, bj.Md5Sum)

	log.Println("Worker", workerId, "processed block", bj.BlockIndex, "with blockId", bj.BlockId)
	log.Printf("%x\n", bj.Md5Sum)

	ldh.Next.Enqueue(bj)
	return nil
}
