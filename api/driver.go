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

// Helper function to quickly convert a byte array to uint16, base 10
// This assumes that the input byte array is a 2byte array (uint16 is 2 bytes)
// Endianess is passed as a second parameter (true for little endian, false for big endian)
func GetUint16FromByte(in []byte, endian bool) uint16 {
	var rs uint16
	if endian {
		// Little Endian
		rs = uint16(in[1])<<8 | uint16(in[0])
	} else {
		// Big Endian
		rs = uint16(in[0])<<8 | uint16(in[1])
	}
	return rs
}

func GetKeyFromButtonMappingByValue(value uint32) (key string, ok bool) {
	for k, v := range ButtonMapping {
		if v == value {
			key = k
			ok = true
			return
		}
	}
	return
}

type Driver struct {
	deviceinfo DeviceInfoInterface
	device     DeviceInterface
	mouse      *IntelliMousePro
}

func (di *Driver) Init(dii DeviceInfoInterface) {
	di.deviceinfo = dii
	di.mouse = &IntelliMousePro{}
}

// Takes control of the device from the kernel.
func (di *Driver) Open() error {
	if di.deviceinfo == nil {
		return errors.New("Device info was not initiated")
	}
	dev, err := di.deviceinfo.Open()
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Unable to open device driver")
	}
	di.device = dev
	return nil
}

// Closes the devuce and returns control to kernel
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

// Debugging function, much easer to just leave this here than have to rewrite it everytime I actually
// need to see something.
func (di *Driver) PrintDebugRead(readHeader uint8) {

	if di.mouse == nil || di.device == nil {
		fmt.Println("Mouse Object or Driver not instantiated")
	}
	fmt.Printf("Using read header %d\n", readHeader)
	// Request for a report from the mouse
	rr_payload, err := di.mouse.TriggerReadRequestPayload(readHeader, []byte{})

	if err != nil {
		fmt.Println("Issue in read payload")
	}
	fmt.Println("Sending payload")
	for _, v := range rr_payload {
		fmt.Printf(fmt.Sprintf("0x%02x ", v))
	}
	fmt.Println("")
	sent, err := di.device.SendFeatureReport(rr_payload)
	if sent != INTELLIMOUSE_PRO_SET_REPORT_LENGTH {
		// TODO: RETURN ERROR.
		fmt.Println("Error when reading report")
	}

	time.Sleep(MOUSE_SLEEP_DRIVER_MILLISECONDS * time.Millisecond)
	// Listen to the report from the mouse
	read_report := make([]byte, INTELLIMOUSE_PRO_GET_REPORT_LENGTH)
	read_report[0] = INTELLIMOUSE_PRO_GET_REPORT
	_, err = di.device.GetFeatureReport(read_report)
	if err != nil {
		fmt.Println("Error when reading report")
	}

	time.Sleep(MOUSE_SLEEP_DRIVER_MILLISECONDS * time.Millisecond)
	fmt.Println(read_report)

	// Cast to hex
	for _, v := range read_report {
		fmt.Printf(fmt.Sprintf("0x%02x ", v))
	}
	fmt.Println("")
}

// SetDpi calls the IntelliMousePro object to calculate an expected payload and send it
// to the HID. It accepts the new dpi value as a parameter.
func (di *Driver) SetDpi(dpi uint16) error {
	if di.mouse == nil || di.device == nil {
		return errors.New("Mouse Object or Device not instantiated")
	}
	dpi_bytearr := di.mouse.SetDpiPayload(dpi)
	write_payload, err := di.mouse.TriggerWriteDataRequestPayload(INTELLIMOUSE_PRO_DPI_WRITE, dpi_bytearr)
	if err != nil {
		return errors.New("Error when creating write report")
	}
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
	rr_payload, err := di.mouse.TriggerReadRequestPayload(INTELLIMOUSE_PRO_DPI_READ, []byte{})

	if err != nil {
		return 0x00, errors.New("Error when reading report")
	}

	sent, err := di.device.SendFeatureReport(rr_payload)
	if sent != INTELLIMOUSE_PRO_SET_REPORT_LENGTH {
		// TODO: RETURN ERROR.
		return 0x00, errors.New("Error when reading report")
	}

	//time.Sleep(MOUSE_SLEEP_DRIVER_MILLISECONDS * time.Millisecond)
	// Listen to the report from the mouse
	read_report := make([]byte, INTELLIMOUSE_PRO_GET_REPORT_LENGTH)
	read_report[0] = INTELLIMOUSE_PRO_GET_REPORT
	_, err = di.device.GetFeatureReport(read_report)
	if err != nil {
		return 0, errors.New("Error when reading report")
	}

	//time.Sleep(MOUSE_SLEEP_DRIVER_MILLISECONDS * time.Millisecond)

	if len(read_report) < 6 {
		return 0x00, errors.New("Result from mouse was too short, unknown error!")
	}
	converted := uint16(read_report[5])<<8 | uint16(read_report[4])

	return converted, nil
}

// ReadLEDColor calls the IntelliMousePro object to get the current LED color on the mouse.
// It returns a six-letter color hex string if successful, else 0
func (di *Driver) ReadLEDColor() (string, error) {
	if di.mouse == nil || di.device == nil {
		return "", errors.New("Mouse Object or Driver not instantiated")
	}
	// Request for a report from the mouse
	rr_payload, err := di.mouse.TriggerReadRequestPayload(INTELLIMOUSE_PRO_COLOR_READ, []byte{})
	if err != nil {
		return "", errors.New("Error when reading report")
	}
	sent, err := di.device.SendFeatureReport(rr_payload)
	if sent != INTELLIMOUSE_PRO_SET_REPORT_LENGTH {
		// TODO: RETURN ERROR.
		return "", errors.New("Error when reading report")
	}

	//time.Sleep(MOUSE_SLEEP_DRIVER_MILLISECONDS * time.Millisecond)
	// Listen to the report from the mouse
	read_report := make([]byte, INTELLIMOUSE_PRO_GET_REPORT_LENGTH)
	read_report[0] = INTELLIMOUSE_PRO_GET_REPORT
	_, err = di.device.GetFeatureReport(read_report)
	if err != nil {
		return "", errors.New("Error when reading report")
	}

	//time.Sleep(MOUSE_SLEEP_DRIVER_MILLISECONDS * time.Millisecond)
	if len(read_report) < 7 {
		return "", errors.New("Read report returned less than 7 bytes")
	}

	hm := HexColor{}
	rs := hm.BigEndianByteToHex(read_report[4:7])
	return rs, nil
}

// SetColor calls the IntelliMousePro object to calculate an expected payload and send it
// to the HID. It accepts the new dpi value as a parameter.
func (di *Driver) SetLEDColor(colorHex string) error {
	if di.mouse == nil || di.device == nil {
		return errors.New("Mouse Object or Device not instantiated")
	}
	color_byte_arr, err := di.mouse.SetColorHexPayload(colorHex)
	if err != nil {
		return err
	}
	write_payload, err := di.mouse.TriggerWriteDataRequestPayload(INTELLIMOUSE_PRO_COLOR_WRITE, color_byte_arr)
	if err != nil {
		return errors.New("Error when creating write report")
	}
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

// Gets the current status of the BackButton on the mouse.
func (di *Driver) GetCurrentBackButton() (string, error) {
	if di.mouse == nil || di.device == nil {
		return "", errors.New("Mouse Object or Device not instantiated")
	}

	get_back_button_arr := di.mouse.GetBackButtonPayload()
	get_back_button_full, err := di.mouse.TriggerReadRequestPayload(INTELLIMOUSE_PRO_BACK_BUTTON_READ, get_back_button_arr)

	if err != nil {
		return "", err
	}

	sent, err := di.device.SendFeatureReport(get_back_button_full)
	if sent != INTELLIMOUSE_PRO_SET_REPORT_LENGTH {
		// TODO: RETURN ERROR.
		return "", errors.New("Error when sending feature report")
	}

	if err != nil {
		return "", err
	}

	//time.Sleep(MOUSE_SLEEP_DRIVER_MILLISECONDS * time.Millisecond)

	read_report := make([]byte, INTELLIMOUSE_PRO_GET_REPORT_LENGTH)
	read_report[0] = INTELLIMOUSE_PRO_GET_REPORT
	_, err = di.device.GetFeatureReport(read_report)

	//time.Sleep(MOUSE_SLEEP_DRIVER_MILLISECONDS * time.Millisecond)

	if len(read_report) < 8 {
		return "", errors.New("Result from mouse was too short, unknown error!")
	}
	converted := uint32(read_report[5])<<24 | uint32(read_report[6])<<16 | uint32(read_report[7])<<8 | uint32(read_report[8])
	stringReprOfButtonMap, ok := GetKeyFromButtonMappingByValue(converted)
	if ok {

		return stringReprOfButtonMap, nil
	}
	return "", errors.New("Unable to find button mapping")
}

// Set back button takes a string representation and performs a lookup to map it to the hexadecimal payload
// An error is returned if it is unable to find the details.
func (di *Driver) SetBackButton(back_button_mapping string) error {
	if di.mouse == nil || di.device == nil {
		return errors.New("Mouse Object or Device not instantiated")
	}
	color_byte_arr, err := di.mouse.SetBackButtonPayload(back_button_mapping)
	if err != nil {
		return err
	}

	write_payload, err := di.mouse.TriggerWriteDataRequestPayload(INTELLIMOUSE_PRO_BACK_BUTTON_WRITE, color_byte_arr)

	if err != nil {
		return errors.New("Error when creating write report")
	}
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
