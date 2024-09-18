package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LocalDataHandler struct {
}

func NewLocalDataHandler() *LocalDataHandler {
	return &LocalDataHandler{}
}

// Process the block
func (ldh *LocalDataHandler) Process(workerId int, bj *Job) error {
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

func (ldh *LocalDataHandler) Start(schedule func(job *Job)) {
	info, err := os.Lstat(config.Path)
	if err != nil {
		log.Println("Error getting file info for", config.Path, ":", err)
		return
	}

	if info.IsDir() {
		log.Println("Error: Path is a directory")
		return
	}

	// Get the file size
	fileSize := info.Size()

	// Calculate the number of blocks
	blockCount := uint32(fileSize / int64(config.BlockSize))
	if fileSize%int64(config.BlockSize) != 0 {
		blockCount++
	}

	if blockCount > 50000 {
		log.Println("Error: File too big")
		return
	}

	// Create a job for each block
	for i := uint32(0); i < blockCount; i++ {
		offset := int64(i) * int64(config.BlockSize)
		job := Job{
			Path:   config.Path,
			Offset: offset,
		}

		schedule(&job)
	}
}
