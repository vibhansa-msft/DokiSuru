package main

import (
	"log"
	"os"
)

type LocalDataHandler struct {
	WorkerPool *WorkerPool
}

func NewLocalDataHandler(wp *WorkerPool) *LocalDataHandler {
	return &LocalDataHandler{
		WorkerPool: wp,
	}
}

// Process the block
func (bj *BlockJob) Process(workerId int) {
	// Open the file
	file, err := os.Open(bj.Path)
	if err != nil {
		log.Println("Worker", workerId, "error opening file", bj.Path, ":", err)
		return
	}
	defer file.Close()

	// Create the block index
	bj.BlockIndex = uint16(bj.Offset / int64(config.BlockSize))

	// Read the data from given offset
	bj.Data = make([]byte, config.BlockSize)
	_, err = file.ReadAt(bj.Data, bj.Offset)
	if err != nil {
		log.Println("Worker", workerId, "error reading file", bj.Path, ":", err)
		return
	}

	// Compuate md5sum of the data
	bj.Md5Sum = ComputeMd5Sum(bj.Data)

	// Convert this slice to a base64 encoded string
	bj.BlockId = GetBlockID(bj.BlockIndex, bj.Md5Sum)
}

func (ldh *LocalDataHandler) Start() {
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
		job := BlockJob{
			Path:   config.Path,
			Offset: offset,
		}

		ldh.WorkerPool.AddJob(&job)
	}
}
