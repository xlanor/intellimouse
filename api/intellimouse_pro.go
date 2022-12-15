package api

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

type IntelliMousePro struct {
	path string `json:"path"`
}

func (imp *IntelliMousePro) Init(path string) {
	imp.path = path
}

func (imp *IntelliMousePro) ToString() string {
	return fmt.Sprintf("%s", imp.path)
}

// Helper function to construct a zeroed out byte array
// for the intellimouse.
func (imp IntelliMousePro) GetSetReportArray() ([]byte, error) {
	if INTELLIMOUSE_PRO_SET_REPORT_LENGTH < 3 {
		// we need a minimum of 3 to ensure that we at least have space for 0, 1, 2
		// aka: report header, property, length of following data.
		// if we can't do that, return an error.
		// guard against acceidentally changing our constants.
		return []byte{}, errors.New("INTELLIMOUSE_PRO_SET_REPORT_LENGTH must be a minimum of 3")
	}
	// byte is just an alias for uint8, which is what 0x00 is (or any numbers that do not exceed 8 bytes)
	request_arr := make([]byte, INTELLIMOUSE_PRO_SET_REPORT_LENGTH)
	for i := 0; i < INTELLIMOUSE_PRO_SET_REPORT_LENGTH; i += 1 {
		request_arr[i] = 0x00
	}
	request_arr[0] = INTELLIMOUSE_PRO_SET_REPORT
	return request_arr, nil
}

// Builds a byte array to send a report to the mouse for a read request.
// Accepts a uint8, the read property that must be specified in the report.
// After this function is called, sleep for 1ms and then try to fetch a report from the HID interface.
func (imp *IntelliMousePro) TriggerReadRequestPayload(read_property uint8, data []byte) ([]byte, error) {
	if len(data) > 0xFF {
		// overflows if we > uint8 max property
		return []byte{}, errors.New("TriggerReadRequestPayload - data is too large")
	}
	request_arr, err := imp.GetSetReportArray()
	if err != nil {
		return []byte{}, err
	}
	request_arr[1] = read_property
	request_arr[2] = uint8(len(data))
	if len(data) > 0 {
		for i := 0; i < len(data); i++ {
			request_arr[i+3] = data[i]
		}
	}

	return request_arr, nil
}

// Builds a byte array to send to the mouse (write).
// Accepts a uint8, the write property that must be specified in the report.
// Accepts a slice of []byte (or uint8 since byte is aliased to that.
// returns a byte slice that has been zeroed out and the respective headers set.
func (imp *IntelliMousePro) TriggerWriteDataRequestPayload(write_property uint8, data []byte) ([]byte, error) {
	if len(data) > 0xFF {
		// overflows if we > uint8 max property
		return []byte{}, errors.New("TriggerWriteDataRequestPayload - data is too large")
	}
	request_arr, err := imp.GetSetReportArray()
	if err != nil {
		return []byte{}, err
	}

	request_arr[1] = write_property
	request_arr[2] = uint8(len(data))
	for i := 0; i < len(data); i++ {
		request_arr[i+3] = data[i]
	}

	return request_arr, nil
}

// Sets the dpi. Clamp this at 200 - 16000
// Perform an implicit clamping instead of explicitly returning an error
// Likewise if it does not match expected steps of a mouse increment in dpi
// This takes up 2 bytes in the packet hence we use uint16
func (imp *IntelliMousePro) SetDpiPayload(dpi uint16) []byte {
	// Make sure that dpi is a multiple
	if dpi%CLAMP_DPI_MULTIPLE != 0 {
		// clamp it to nearest multiple
		dpi = dpi + (CLAMP_DPI_MULTIPLE - dpi%CLAMP_DPI_MULTIPLE)
	}

	if dpi > 16000 {
		dpi = 16000
	}
	if dpi < 200 {
		dpi = 200
	}
	// Convert to bytes
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, dpi)
	return bs
}

// Gets the data payload for setting a color hex LED after validating it.
func (imp *IntelliMousePro) SetColorHexPayload(colorHexCode string) ([]byte, error) {
	h := HexColor{}
	h.Init(colorHexCode)
	parsed, err := h.ValidateHex()
	// this is naturally in big endian which is fine
	if err != nil {
		return []byte{}, err
	}
	data, err := hex.DecodeString(parsed)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

// Constructs a payload to send to the mouse to find out what the current back button is
func (imp *IntelliMousePro) GetBackButtonPayload() []byte {
	return []byte{INTELLIMOUSE_PRO_BACK_BUTTON}
}

// Accepts a button map string (since we expose this in the UI) and checks against our constant map.
func (imp *IntelliMousePro) SetBackButtonPayload(buttonMap string) ([]byte, error) {
	ret := make([]byte, 4)
	if val, ok := ButtonMapping[buttonMap]; ok {
		binary.BigEndian.PutUint32(ret, val)
		return append([]byte{INTELLIMOUSE_PRO_BACK_BUTTON}, ret...), nil
	} else {
		return []byte{}, errors.New("Cannot find a valid button mapping")
	}
}

// Constructs a payload to send to the mouse to find out what the current middle button is
func (imp *IntelliMousePro) GetMiddleButtonPayload() []byte {
	return []byte{INTELLIMOUSE_PRO_MIDDLE_BUTTON}
}

// Accepts a button map string (since we expose this in the UI) and checks against our constant map.
func (imp *IntelliMousePro) SetMiddleButtonPayload(buttonMap string) ([]byte, error) {
	ret := make([]byte, 4)
	if val, ok := ButtonMapping[buttonMap]; ok {
		binary.BigEndian.PutUint32(ret, val)
		return append([]byte{INTELLIMOUSE_PRO_MIDDLE_BUTTON}, ret...), nil
	} else {
		return []byte{}, errors.New("Cannot find a valid button mapping")
	}
}
