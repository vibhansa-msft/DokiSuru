package main

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

func getServiceClient(account string) (*service.Client, error) {
	key := strings.TrimSpace(os.Getenv("AZURE_STORAGE_KEY"))
	sas := strings.TrimSpace(os.Getenv("AZURE_STORAGE_SAS"))

	if key == "" && sas == "" {
		log.Println("Either access key or SAS is needed")
		return nil, fmt.Errorf("Either access key or SAS is needed")
	}

	var svcClient *service.Client
	var err error
	svcURL := "https://" + account + ".blob.core.windows.net/"
	if key != "" {
		cred, err := service.NewSharedKeyCredential(account, key)
		if err != nil {
			log.Printf("Unable to create shared key [%v]", err.Error())
			return nil, err
		}

		svcClient, err = service.NewClientWithSharedKeyCredential(svcURL, cred, nil)
		if err != nil {
			log.Printf("Unable to create service client [%v]", err.Error())
			return nil, err
		}
	} else {
		svcURL = svcURL + "?" + sas
		svcClient, err = service.NewClientWithNoCredential(svcURL, nil)
		if err != nil {
			log.Printf("Unable to create service client [%v]", err.Error())
			return nil, err
		}
	}

	return svcClient, nil
}

func GetContainerClient() (*container.Client, error) {
	account := strings.TrimSpace(os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"))
	container := strings.TrimSpace(os.Getenv("AZURE_STORAGE_ACCOUNT_CONTAINER"))

	if account == "" || container == "" {
		log.Println("Storage account name or container not provided")
		return nil, fmt.Errorf("Storage account name or container not provided")
	}

	svcClient, err := getServiceClient(account)
	if err != nil {
		return nil, err
	}

	return svcClient.NewContainerClient(container), err
}
