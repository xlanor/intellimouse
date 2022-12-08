package main

import (
	"intellimouse/api"
)

// this is temporary and we will get rid of this and put wails in.
func main() {
	devices := api.GetDeviceInfo(api.HidStruct{})
	if len(devices) < 1 {

	} else {
		dri := api.Driver{}
		dri.Init(devices[0])
	}
}
