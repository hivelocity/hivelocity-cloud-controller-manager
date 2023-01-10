package client

import (
	hv "github.com/hivelocity/hivelocity-client-go/client"
)

type API interface {
	GetBareMetalDeviceIdResource(deviceId int32) (*hv.BareMetalDevice, error)
}
