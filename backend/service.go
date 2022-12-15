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
	a.Log.Info(fmt.Sprintf("Received checksum %s", checksum))
	if deviceInfo, ok := a.AvaliableDevices[checksum]; ok {
		if deviceInfo.DeviceInfo != nil {
			a.Log.Info("Opening driver")
			a.Driver = &api.Driver{}
			a.Driver.Init(*deviceInfo.DeviceInfo)
			err := a.Driver.Open()
			if err != nil {
				return err
			} else {
				a.Log.Info("Driver successfully opened")
				a.GetDeviceInformation()
			}
		} else {
			a.Log.Error("Device Information has no avaliable devices")
			return errors.New("Unable to Open device")
		}

	} else {
		a.Log.Error("Device array does not have any avaliable devices")
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
		m := MouseInformationStruct{}
		m.Init(a.Driver)
		a.Log.Info(fmt.Sprintf("DPI: %d", m.Dpi))
		a.Log.Info(fmt.Sprintf("LED: %s", m.LedColor))
		a.Log.Info("Emitting mouse information to frontend")
		runtime.EventsEmit(a.ctx, "mouseinformation", m)
	} else {
		a.Log.Error("Context Driver is not set")
	}
}

func (a *App) SetLEDWrapper(ledhex string) {
	if a.Driver != nil {
		err := a.Driver.SetLEDColor(ledhex)
		if err != nil {
			a.Log.Error(fmt.Sprintf("Error when setting LED: %s", err.Error()))
		} else {
			a.Log.Info(fmt.Sprintf("LED changed to %s\n", ledhex))
		}
	} else {
		a.Log.Error("Context Driver is not set")
	}
}

func (a *App) SetDpiWrapper(dpi int) {
	dpi16 := uint16(dpi)
	if a.Driver != nil {
		err := a.Driver.SetDpi(dpi16)
		if err != nil {
			a.Log.Error(fmt.Sprintf("Error when setting DPI: %s", err.Error()))
		} else {
			a.Log.Info(fmt.Sprintf("DPI changed to %d", dpi16))
		}
	} else {
		a.Log.Error("Context Driver is not set")
	}
}

func (a *App) SetButtonWrapper(button_type string, new_value string) {
	if a.Driver != nil {
		if button_type == "back" {
			err := a.Driver.SetBackButton(api.BACK_BUTTON, new_value)
			if err != nil {
				a.Log.Error(err.Error())
			}else {
				a.Log.Info(fmt.Sprintf("Button set to %s", new_value))
			}
		}else if button_type == "front" {

		} else if button_type == "middle"{ 

		} else {
			a.Log.Error(fmt.Sprintf("Unknown button type: Received %s", button_type))
		}
	} else {
		a.Log.Error("Context driver is not set")
	}
}