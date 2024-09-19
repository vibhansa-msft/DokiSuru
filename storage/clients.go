package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

// environment variables:
// AZURE_STORAGE_ACCOUNT_NAME
// AZURE_STORAGE_ACCOUNT_CONTAINER
// AZURE_STORAGE_KEY
// AZURE_STORAGE_SAS

type Clients struct {
	svcClient *service.Client
	cntClient *container.Client
}

type StgBlock struct {
	Name string
	Size int64
}

func NewClients() (c *Clients, err error) {
	c = &Clients{}
	err = c.createServiceClient()
	if err != nil {
		return nil, err
	}
	err = c.createContainerClient()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Clients) GetContainerClient() *container.Client {
	return c.cntClient
}

func (c *Clients) createServiceClient() error {
	account := strings.TrimSpace(os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"))
	container := strings.TrimSpace(os.Getenv("AZURE_STORAGE_ACCOUNT_CONTAINER"))

	if account == "" || container == "" {
		log.Println("Storage account name or container not provided")
		return fmt.Errorf("Storage account name or container not provided")
	}

	key := strings.TrimSpace(os.Getenv("AZURE_STORAGE_KEY"))
	sas := strings.TrimSpace(os.Getenv("AZURE_STORAGE_SAS"))

	if key == "" && sas == "" {
		log.Println("Either access key or SAS is needed")
		return fmt.Errorf("Either access key or SAS is needed")
	}

	var svcClient *service.Client
	var err error
	svcURL := "https://" + account + ".blob.core.windows.net/"
	if key != "" {
		cred, err := service.NewSharedKeyCredential(account, key)
		if err != nil {
			log.Printf("Unable to create shared key [%v]", err.Error())
			return err
		}

		svcClient, err = service.NewClientWithSharedKeyCredential(svcURL, cred, nil)
		if err != nil {
			log.Printf("Unable to create service client [%v]", err.Error())
			return err
		}
	} else {
		svcURL = svcURL + "?" + sas
		svcClient, err = service.NewClientWithNoCredential(svcURL, nil)
		if err != nil {
			log.Printf("Unable to create service client [%v]", err.Error())
			return err
		}
	}
	c.svcClient = svcClient
	return nil
}

func (c *Clients) createContainerClient() error {
	container := strings.TrimSpace(os.Getenv("AZURE_STORAGE_ACCOUNT_CONTAINER"))
	if container == "" {
		log.Println("Storage container not provided")
		return fmt.Errorf("Storage container not provided")
	}
	c.cntClient = c.svcClient.NewContainerClient(container)

	return nil
}

func (c *Clients) CreateBlobClient(name string) *blockblob.Client {
	return c.cntClient.NewBlockBlobClient(name)
}

func (c *Clients) GetBlockList(bbc *blockblob.Client) ([]StgBlock, error) {
	prop, err := bbc.GetProperties(context.Background(), nil)
	if err != nil {
		return []StgBlock{}, nil
	}

	val, ok := prop.Metadata["Doki_etag"]
	if !ok || (string)(*prop.ETag) == *val {
		return nil, fmt.Errorf("ETag mismatch")
	}

	resp, err := bbc.GetBlockList(context.Background(), blockblob.BlockListTypeCommitted, nil)
	if err != nil {
		log.Println("Error getting block list:", err)
		return nil, err
	}

	var blocks []StgBlock
	for _, block := range resp.CommittedBlocks {
		blocks = append(blocks, StgBlock{Name: *block.Name, Size: *block.Size})
	}

	return blocks, nil
}

func (c *Clients) StageBlock(bbc *blockblob.Client, data []byte, id string) error {
	_, err := bbc.StageBlock(context.Background(), id, streaming.NopCloser(bytes.NewReader(data)), nil)
	return err
}

func (c *Clients) GetBlock(bbc *blockblob.Client, offset int64, length int64) ([]byte, error) {
	data := make([]byte, length)
	_, err := bbc.DownloadBuffer(context.Background(), data, &blob.DownloadBufferOptions{
		Range: blob.HTTPRange{
			Offset: offset,
			Count:  length,
		}})
	return data, err
}

func (c *Clients) PutBlockList(bbc *blockblob.Client, list []string) error {
	_, err := bbc.CommitBlockList(context.Background(), list, nil)
	if err != nil {
		log.Println("Error committing block list:", err)
		return err
	}

	resp, err := bbc.GetProperties(context.Background(), nil)
	if err != nil {
		return err
	}

	resp.Metadata["Doki_etag"] = (*string)(resp.ETag)
	_, err = bbc.SetMetadata(context.Background(), resp.Metadata, nil)

	return err
}
