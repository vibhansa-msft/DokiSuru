package main

import (
	"dokisuru/storage"
	"dokisuru/utils"
	"log"
	"os"
	"strings"

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

		id_str := strings.Join(list, ",")
		id_hash := utils.ComputeMd5Sum([]byte(id_str))

		err := rdh.Client.PutBlockList(rdh.Blob.Client, list, id_hash)
		if err != nil {
			log.Println("Error committing block list")
			return err
		}
	} else {
		log.Println("No changes to commit")
	}

	if config.Validate {
		rdh.validate()
	}

	return nil
}

// Process the block
func (rdh *RemoteDataHandler) Process(workerId int, bj *Job) error {
	if bj.BlockIndex == 0 && len(rdh.Blob.Blocks) != int(bj.NoOfBlocks) {
		rdh.Blob.Modified = true
		rdh.Blob.Blocks = rdh.Blob.Blocks[:bj.NoOfBlocks]
	}

	if bj.BlockId == "" {
		return nil
	}

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

// Process the block
func (rdh *RemoteDataHandler) validate() {
	log.Println("Validating blob")
	err := rdh.Client.DownloadBlob(rdh.Blob.Client, "_validate_blob.bin")
	if err != nil {
		log.Println("Error downloading blob:", err)
		return
	}
	log.Println("Download complete")

	remote_md5 := utils.GetMd5File("_validate_blob.bin")
	local_md5 := utils.GetMd5File(config.Path)

	log.Println("Remote MD5:", remote_md5, "Local MD5:", local_md5)

	if remote_md5 != local_md5 {
		log.Println("Blob mismatch")
	} else {
		log.Println("Blob validated successfully")
	}

	os.Remove("_validate_blob.bin")
}
