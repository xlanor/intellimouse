package backend

import (
	"fmt"
	"github.com/dolmen-go/hid"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xlanor/intellimouse/api"
	"golang.org/x/exp/maps"
	"time"
)

func inslice(n uint64, h []uint64) bool {
	for _, v := range h {
		if v == n {
			return true
		}
	}
	return false
}

// expose to vue
type DeviceInformationJson struct {
	Path         string          `json:"path"`
	VendorID     uint16          `json:"vendor_id"`
	ProductID    uint16          `json:"product_id"`
	Serial       string          `json:"serial"`
	Manufacturer string          `json:"manufacturer"`
	Product      string          `json:"product"`
	Interface    int             `json:"interface"`
	Hash         uint64          `json:"checksum"`
	DeviceInfo   *hid.DeviceInfo // Keep this in memory first to open a device easily
}

func (d *DeviceInformationJson) ParseFromHidLib(di *hid.DeviceInfo) {
	d.Path = di.Path
	d.VendorID = di.VendorID
	d.ProductID = di.ProductID
	d.Serial = di.Serial
	d.Manufacturer = di.Manufacturer
	d.Product = di.Product
	d.Interface = di.Interface
	d.DeviceInfo = nil
	hash, err := hashstructure.Hash(d, hashstructure.FormatV2, nil)
	if err != nil {
		panic(fmt.Sprintf("ParseFromHdbLib error: %s", err.Error()))
	}
	d.Hash = hash
	d.DeviceInfo = di
}

func (a *App) UpdateAvaliableDevices(dij []DeviceInformationJson) {
	// build an array of hashes
	hashes := make([]uint64, 0)
	for _, v := range dij {
		if _, ok := a.AvaliableDevices[v.Hash]; !ok {
			a.AvaliableDevices[v.Hash] = v
		}
		hashes = append(hashes, v.Hash)
	}
	map_keys := maps.Keys(a.AvaliableDevices)
	res := make([]uint64, 0)
	for _, v := range map_keys {
		if !inslice(v, hashes) {
			res = append(res, v)
		}
	}
	fmt.Println("DELETE RES %v", res)
	for _, toDelete := range res {
		delete(a.AvaliableDevices, toDelete)
	}
	return
}

func (a *App) LoadAvaliableDevices() []DeviceInformationJson {
	ret := make([]DeviceInformationJson, 0)
	avaliableDevices := api.GetDeviceInfo(api.HidStruct{})
	if len(avaliableDevices) == 0 {
		a.UpdateAvaliableDevices(ret)
		return ret
	}
	for _, v := range avaliableDevices {
		dij := DeviceInformationJson{}
		dij.ParseFromHidLib(&v)
		ret = append(ret, dij)
	}
	a.UpdateAvaliableDevices(ret)
	return ret
}

func (a *App) LoadDevicesPolling() error {
	go func() {
		for {
			fmt.Println("Polling HID interface")
			runtime.EventsEmit(a.ctx, "devices", a.LoadAvaliableDevices())
			time.Sleep(1 * time.Second)
			fmt.Printf("%v\n", a.AvaliableDevices)
		}
	}()
	return nil
}
