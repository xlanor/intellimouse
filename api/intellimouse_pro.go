package api

import (
	"encoding/binary"
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

// Builds a byte array to send a report to the mouse for a read request.
// Accepts a uint8, the read property that must be specified in the report.
// After this function is called, sleep for 1ms and then try to fetch a report from the HID interface.
func (imp *IntelliMousePro) TriggerReadRequestPayload(read_property uint8) []byte {
	// byte is just an alias for uint8, which is what 0x00 is (or any numbers that do not exceed 8 bytes)
	request_arr := make([]byte, INTELLIMOUSE_PRO_SET_REPORT_LENGTH)
	for i := 0; i < INTELLIMOUSE_PRO_SET_REPORT_LENGTH; i += 1 {
		request_arr[i] = 0x00
	}
	request_arr[0] = INTELLIMOUSE_PRO_SET_REPORT
	request_arr[1] = read_property
	request_arr[2] = 0x01

	return request_arr
}

func (imp *IntelliMousePro) TriggerWriteDataRequestPayload(write_property uint8, data []byte) []byte {
	request_arr := make([]byte, INTELLIMOUSE_PRO_SET_REPORT_LENGTH)
	for i := 0; i < INTELLIMOUSE_PRO_SET_REPORT_LENGTH; i += 1 {
		request_arr[i] = 0x00
	}
	request_arr[0] = INTELLIMOUSE_PRO_SET_REPORT
	request_arr[1] = write_property
	request_arr[2] = uint8(len(data))
	for i := 0; i < len(data); i++ {
		request_arr[i+3] = data[i]
	}
	return request_arr
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

func (imp *IntelliMousePro) GetDpi() []byte {

	return []byte{}
}
