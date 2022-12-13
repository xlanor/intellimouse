package backend

import (
	"errors"
	"fmt"
	"github.com/dolmen-go/hid"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xlanor/intellimouse/api"
	"golang.org/x/exp/maps"
	"strconv"
	"time"
)

func inslice(n string, h []string) bool {
	for _, v := range h {
		if v == n {
			return true
		}
	}
	return false
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
	} else {

	}
	// dont use uint, js will overflow it...
	d.Hash = strconv.FormatUint(hash, 10)
	d.DeviceInfo = di
}

func (a *App) UpdateAvaliableDevices(dij []DeviceInformationJson) {
	// build an array of hashes
	hashes := make([]string, 0)
	for _, v := range dij {
		if _, ok := a.AvaliableDevices[v.Hash]; !ok {
			a.AvaliableDevices[v.Hash] = v
		}
		hashes = append(hashes, v.Hash)
	}
	map_keys := maps.Keys(a.AvaliableDevices)
	res := make([]string, 0)
	for _, v := range map_keys {
		if !inslice(v, hashes) {
			res = append(res, v)
		}
	}
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

func (a *App) SelectDevice(checksum string) error {
	fmt.Println("Checksum here")
	fmt.Println(checksum)
	if deviceInfo, ok := a.AvaliableDevices[checksum]; ok {
		if deviceInfo.DeviceInfo != nil {
			fmt.Println("Opening driver")
			a.Driver = &api.Driver{}
			a.Driver.Init(*deviceInfo.DeviceInfo)
			err := a.Driver.Open()
			if err != nil {
				return err
			} else {
				fmt.Println(`Driver opened`)
				a.GetDeviceInformation()
			}
		} else {
			fmt.Println("Device Info has no avaliable device")
			fmt.Printf("%v\n", deviceInfo)
			return errors.New("Unable to Open device")
		}

	} else {
		fmt.Println("Device Array has no avaliable device")
		fmt.Printf("%v\n", a.AvaliableDevices)
		return errors.New("Unable to find device")
	}
	return nil
}

func (a *App) LoadDevicesPolling() error {
	// Poll for new usb devices
	go func() {
		for {
			runtime.EventsEmit(a.ctx, "devices", a.LoadAvaliableDevices())
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}

func (a *App) GetDeviceInformation() {
	if a.Driver != nil {
		fmt.Println("Here")
		m := MouseInformationStruct{}
		m.Init(a.Driver)
		fmt.Printf("%v\n", m)
		runtime.EventsEmit(a.ctx, "mouseinformation", m)
	}else {
		fmt.Println("Driver nil")
	}
}