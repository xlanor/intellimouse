package backend

import (
	"fmt"
	"github.com/xlanor/intellimouse/api"
	"github.com/dolmen-go/hid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"time"
)

// expose to vue
type DeviceInformationJson struct {
	Path string `json:"path"`
	VendorID uint16 `json:"vendor_id"`
	ProductID uint16 `json:"product_id"`
	Serial		string `json:"serial"`
	Manufacturer string `json:"manufacturer"`
	Product string `json:"product"`
	Interface int `json:"interface"`
	DeviceInfo *hid.DeviceInfo	// Keep this in memory first to open a device easily
}

func (d *DeviceInformationJson) ParseFromHidLib(di *hid.DeviceInfo) {
	d.Path = di.Path
	d.VendorID = di.VendorID
	d.ProductID = di.ProductID
	d.Serial = di.Serial
	d.Manufacturer = di.Manufacturer
	d.Product = di.Product
	d.Interface = di.Interface
	d.DeviceInfo = di
}

func (a *App) LoadAvaliableDevices() ([]DeviceInformationJson) {
	ret := make([]DeviceInformationJson, 0)
	avaliableDevices := api.GetDeviceInfo(api.HidStruct{})
	if len(avaliableDevices) == 0 {
		return ret
	}
	for _, v := range avaliableDevices {
		dij := DeviceInformationJson{}
		dij.ParseFromHidLib(&v)
		ret = append(ret, dij)
	}
	return ret

}

func (a *App) LoadDevicesPolling() error {
	go func() {
		for {
			fmt.Println("Polling HID interface")
			runtime.EventsEmit(a.ctx, "devices", a.LoadAvaliableDevices())
			time.Sleep(15 * time.Second)
		}
	}()
	return nil
}