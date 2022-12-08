package api

import (
	"errors"
	"fmt"
	"github.com/dolmen-go/hid"
	"time"
)

// GetDeviceInfo accepts a HidInterface object and returns an array of dolmen-go/hid.DeviceInfo
// If there is no matching interface found, an empty array will be returned.
func GetDeviceInfo(hidInterface HidInterface) []hid.DeviceInfo {
	rs := make([]hid.DeviceInfo, 0)
	di := hidInterface.Enumerate(INTELLIMOUSE_PRO_VENDOR_ID, INTELLIMOUSE_PRO_PRODUCT_ID)
	for _, v := range di {
		if v.Interface == INTELLIMOUSE_PRO_INTERFACE {
			rs = append(rs, v)
		}
	}
	return rs
}

type Driver struct {
	deviceinfo DeviceInfoInterface
	device     DeviceInterface
	mouse      *IntelliMousePro
}

// Left as non pointer, but it accepts a Pointer as an input too.
func (di *Driver) Init(dii DeviceInfoInterface) {
	di.deviceinfo = dii
	di.mouse = &IntelliMousePro{}
}

func (di *Driver) Open() error {
	if di.deviceinfo == nil {
		return errors.New("Device info was not initiated")
	}
	dev, err := di.deviceinfo.Open()
	if err != nil {
		return errors.New("Unable to open device driver")
	}
	di.device = dev
	return nil
}

func (di *Driver) Close() error {
	if di.device == nil {
		return errors.New("No device currently present")
	}
	err := di.device.Close()
	if err != nil {
		return err
	}
	return nil
}

// SetDpi calls the IntelliMousePro object to calculate an expected payload and send it
// to the HID. It accepts the new dpi value as a parameter.
func (di *Driver) SetDpi(dpi uint16) error {
	if di.mouse == nil || di.device == nil {
		return errors.New("Mouse Object or Device not instantiated")
	}
	dpi_bytearr := di.mouse.SetDpiPayload(dpi)
	write_payload := di.mouse.TriggerWriteDataRequestPayload(INTELLIMOUSE_PRO_DPI_WRITE, dpi_bytearr)
	sent, err := di.device.SendFeatureReport(write_payload)
	if sent != INTELLIMOUSE_PRO_SET_REPORT_LENGTH {
		return errors.New(fmt.Sprintf("sent bytes did not match write payload, sent:%d expected:%d", sent, len(write_payload)))
	}
	if err != nil {
		return err
	}
	time.Sleep(MOUSE_SLEEP_DRIVER_MILLISECONDS * time.Millisecond)
	return nil
}

// ReadDpi calls the IntelliMousePro object to get the current dpi on the object.
// It returns a uint16 if successful, else 0
func (di *Driver) ReadDpi() (uint16, error) {
	if di.mouse == nil || di.device == nil {
		return 0, errors.New("Mouse Object or Driver not instantiated")
	}
	// Request for a report from the mouse
	rr_payload := di.mouse.TriggerReadRequestPayload(INTELLIMOUSE_PRO_DPI_READ)
	sent, err := di.device.SendFeatureReport(rr_payload)
	if sent != INTELLIMOUSE_PRO_SET_REPORT_LENGTH {
		// TODO: RETURN ERROR.
		return 0, errors.New("Error when reading report")
	}

	time.Sleep(MOUSE_SLEEP_DRIVER_MILLISECONDS * time.Millisecond)
	// Listen to the report from the mouse
	read_report := make([]byte, INTELLIMOUSE_PRO_GET_REPORT_LENGTH)
	read_report[0] = INTELLIMOUSE_PRO_GET_REPORT
	_, err = di.device.GetFeatureReport(read_report)
	if err != nil {
		return 0, errors.New("Error when reading report")
	}

	time.Sleep(MOUSE_SLEEP_DRIVER_MILLISECONDS * time.Millisecond)

	if len(read_report) < 6 {
		return 0, errors.New("Result from mouse was too short, unknown error!")
	}
	// TODO: Check whether the array length is right before accessing unsafely
	converted := uint16(read_report[5])<<8 | uint16(read_report[4])

	return converted, nil
}
