package main

import (
	"fmt"
	"github.com/xlanor/intellimouse/api"
	"os"
	"time"
)

// for me to quickly test out stuff.
func main() {
	devices := api.GetDeviceInfo(api.HidStruct{})
	if len(devices) < 1 {
		fmt.Println("Could not find a valid mouse")
		os.Exit(1)
	}
	device := devices[0]
	driver := api.Driver{}
	driver.Init(device)
	err := driver.Open()
	if err != nil {
		fmt.Println("Could not load driver")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer driver.Close()
	fmt.Println(driver.ReadDpi())
	fmt.Println(driver.ReadLEDColor())
	cyan := "#d40b58"
	time.Sleep(1000 * time.Millisecond)
	err = driver.SetLEDColor(cyan)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(driver.GetCurrentBackButton())
	driver.SetBackButton("INTELLIMOUSE_PRO_BACK_BUTTON_SET_BACK_BUTTON")
	fmt.Println("Button set")
}
