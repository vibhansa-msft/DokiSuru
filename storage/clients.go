package storage

import (
	"fmt"
	"log"
	"os"
	"strings"

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
