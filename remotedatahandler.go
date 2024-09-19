package main

import (
	"dokisuru/storage"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
)

var _ JobHandler = &RemoteDataHandler{}

type BlobDetails struct {
	Client   *blockblob.Client
	Name     string
	Blocks   []storage.StgBlock
	Modified bool
}

type RemoteDataHandler struct {
	BaseHandler
	Blob   BlobDetails
	Client *storage.Clients
}

func NewRemoteDataHandler(workerCount int, next JobHandler) *RemoteDataHandler {
	rdh := &RemoteDataHandler{
		BaseHandler: BaseHandler{
			Next: next},
	}
	rdh.Worker = NewWorkerPool(workerCount, rdh.Process)
	rdh.Worker.Start()
	return rdh
}

func (rdh *RemoteDataHandler) Start() error {
	var err error
	rdh.Client, err = storage.NewClients()
	if err != nil {
		log.Println("Error creating clients:", err)
		return err
	}

	rdh.Blob = BlobDetails{
		Client:   rdh.Client.CreateBlobClient(config.Path),
		Name:     config.Path,
		Modified: false,
	}

	rdh.Blob.Blocks, err = rdh.Client.GetBlockList(rdh.Blob.Client)
	if err != nil {
		log.Println("Error getting block list")
		return err
	}

	if len(rdh.Blob.Blocks) == 0 {
		rdh.Blob.Blocks = make([]storage.StgBlock, 50000)
	}

	return nil
}

func (rdh *RemoteDataHandler) Stop() error {
	rdh.Worker.Stop()

	if rdh.Blob.Modified {
		list := make([]string, 0)
		for _, block := range rdh.Blob.Blocks {
			if block.Name == "" {
				break
			}
			list = append(list, block.Name)
		}

		err := rdh.Client.PutBlockList(rdh.Blob.Client, list)
		if err != nil {
			log.Println("Error committing block list")
			return err
		}
	}

	return nil
}

// Process the block
func (rdh *RemoteDataHandler) Process(workerId int, bj *Job) error {
	if int(bj.BlockIndex) < len(rdh.Blob.Blocks) {
		if rdh.Blob.Blocks[bj.BlockIndex].Name == bj.BlockId {
			return nil
		}
	}

	err := rdh.Client.StageBlock(rdh.Blob.Client, bj.Data, bj.BlockId)
	if err != nil {
		log.Println("Worker", workerId, "error staging block", bj.BlockIndex, ":", err)
		return err
	}

	log.Println("Worker", workerId, "processed block", bj.BlockIndex, "with blockId", bj.BlockId)
	log.Printf("%x\n", bj.Md5Sum)

	rdh.Blob.Blocks[bj.BlockIndex].Name = bj.BlockId
	rdh.Blob.Modified = true
	return nil
}
