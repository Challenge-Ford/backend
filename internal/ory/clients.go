package ory

import (
	ketoclient "github.com/ory/keto-client-go"
	kratosclient "github.com/ory/kratos-client-go"
)

const (
	DealershipsNS = "dealerships"
	PermissionsNS = "permissions"
)

func NewKratosClient(adminURL string) *kratosclient.APIClient {
	cfg := kratosclient.NewConfiguration()
	cfg.Servers = kratosclient.ServerConfigurations{{URL: adminURL}}
	return kratosclient.NewAPIClient(cfg)
}

func NewKetoReadClient(readURL string) *ketoclient.APIClient {
	cfg := ketoclient.NewConfiguration()
	cfg.Servers = ketoclient.ServerConfigurations{{URL: readURL}}
	return ketoclient.NewAPIClient(cfg)
}

func NewKetoWriteClient(writeURL string) *ketoclient.APIClient {
	cfg := ketoclient.NewConfiguration()
	cfg.Servers = ketoclient.ServerConfigurations{{URL: writeURL}}
	return ketoclient.NewAPIClient(cfg)
}
