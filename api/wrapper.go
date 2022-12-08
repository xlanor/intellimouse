package api

import (
	"github.com/dolmen-go/hid"
)

// Use all references to dolmen-go/hid as an interface to allow for mocking via struct substitution
type HidInterface interface {
	Enumerate(vendorID uint16, productID uint16) []hid.DeviceInfo
}

// HidStruct wraps dolmen-go/hid to expose the functions for mocking.
type HidStruct struct{}

func (hs HidStruct) Enumerate(vendorID uint16, productID uint16) []hid.DeviceInfo {
	return hid.Enumerate(vendorID, productID)
}

// Allows us to mock a hid.DeviceInfo
type DeviceInfoInterface interface {
	Open() (*hid.Device, error)
}

// Allows us to mock a *hid.Device
// As the Device struct is mocked with pointer receivers this means that
// all usage of this interface should be utilised with a pointer value
type DeviceInterface interface {
	Close() error
	GetFeatureReport(b []byte) (int, error)
	Read(b []byte) (int, error)
	SendFeatureReport(b []byte) (int, error)
	Write(b []byte) (int, error)
}
