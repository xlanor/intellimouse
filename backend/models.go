package backend

import (
	"fmt"
	"github.com/dolmen-go/hid"
	"github.com/xlanor/intellimouse/api"
)

// exposes information to vue

type DeviceInformationJson struct {
	Path         string          `json:"path"`
	VendorID     uint16          `json:"vendor_id"`
	ProductID    uint16          `json:"product_id"`
	Serial       string          `json:"serial"`
	Manufacturer string          `json:"manufacturer"`
	Product      string          `json:"product"`
	Interface    int             `json:"interface"`
	Hash         string          `json:"checksum"`
	DeviceInfo   *hid.DeviceInfo // Keep this in memory first to open a device easily
}

type MouseInformationStruct struct {
	Dpi               uint16 `json:"dpi"`
	CurrentBackButton string `json:"back_button"`
	LedColor          string `json:"led"`
}

func (m *MouseInformationStruct) Init(d *api.Driver) {
	if d != nil {
		dpi, err := d.ReadDpi()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			m.Dpi = dpi
		}
		back, err := d.GetButtonPayload(api.BACK_BUTTON)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			m.CurrentBackButton = back
		}
		led, err := d.ReadLEDColor()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			m.LedColor = led
		}
	}
}
