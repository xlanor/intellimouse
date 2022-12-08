package api

import (
	"errors"
	"github.com/dolmen-go/hid"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	VALID_PATH   = "VALID_PATH"
	INVALID_PATH = "INVALID_PATH"
)

type HidMockStruct struct {
	set_non_valid_enumerate bool
}

func (hms *HidMockStruct) Init(mockEnumerate bool) {
	hms.set_non_valid_enumerate = mockEnumerate
}

func (hms HidMockStruct) Enumerate(vendorID uint16, productID uint16) []hid.DeviceInfo {
	rs := make([]hid.DeviceInfo, 0)
	if hms.set_non_valid_enumerate {
		rs = append(rs, hid.DeviceInfo{
			Path:      VALID_PATH,
			Interface: INTELLIMOUSE_PRO_INTERFACE,
		})
	}
	// add invalid enumerates
	for i := 0; i < 10; i++ {
		rs = append(rs, hid.DeviceInfo{
			Path:      INVALID_PATH,
			Interface: INTELLIMOUSE_PRO_INTERFACE + 0x01,
		})
	}
	return rs
}

// Test using the hid driver to find a valid mouse interface
func Test_Driver_GetDeviceInfo(t *testing.T) {
	t.Parallel()
	t.Run("Contains one valid device", func(t *testing.T) {
		// mock an interface.
		hms := HidMockStruct{}
		hms.Init(true)
		rs := GetDeviceInfo(hms)
		assert.Equal(t, 1, len(rs))
		assert.Equal(t, VALID_PATH, rs[0].Path)
	})
	t.Run("Contains no valid device", func(t *testing.T) {
		// mock an interface.
		hms := HidMockStruct{}
		hms.Init(false)
		rs := GetDeviceInfo(hms)
		assert.Equal(t, 0, len(rs))
	})
}

type DeviceInfoStructWithOpenError struct{}

func (dis DeviceInfoStructWithOpenError) Open() (*hid.Device, error) {
	return nil, errors.New("device open error")
}

type DeviceInfoStructWithOpenOk struct{}

func (dis DeviceInfoStructWithOpenOk) Open() (*hid.Device, error) { return nil, nil }

type OpenDeviceStructWithCloseError struct{}

func (ods OpenDeviceStructWithCloseError) Close() error {
	return errors.New("device close error")
}
func (ods OpenDeviceStructWithCloseError) GetFeatureReport(b []byte) (int, error)  { return 0, nil }
func (ods OpenDeviceStructWithCloseError) Read(b []byte) (int, error)              { return 0, nil }
func (ods OpenDeviceStructWithCloseError) SendFeatureReport(b []byte) (int, error) { return 0, nil }
func (ods OpenDeviceStructWithCloseError) Write(b []byte) (int, error)             { return 0, nil }

func Test_Driver_Open(t *testing.T) {
	t.Parallel()
	t.Run("Open device driver failed - not set", func(t *testing.T) {
		d := Driver{}
		err := d.Open()
		assert.Error(t, err)
		assert.Equal(t, "Device info was not initiated", err.Error())
	})
	t.Run("Open device driver failed - device busy or otherwise", func(t *testing.T) {
		d := Driver{}
		di := DeviceInfoStructWithOpenError{}
		d.Init(di)
		err := d.Open()
		assert.Error(t, err)
		assert.Equal(t, "Unable to open device driver", err.Error())
	})
	t.Run("Open device driver ok", func(t *testing.T) {
		d := Driver{}
		d.Init(DeviceInfoStructWithOpenOk{})
		err := d.Open()
		assert.Nil(t, err)
	})
}

func Test_Driver_Close(t *testing.T) {
	t.Parallel()
	t.Run("Open device driver failed - not set", func(t *testing.T) {
		d := Driver{}
		err := d.Close()
		assert.Error(t, err)
		assert.Equal(t, "No device currently present", err.Error())
	})
}

type SetDpiDeviceStructWithSetReportError struct{}

func (ods SetDpiDeviceStructWithSetReportError) Close() error { return nil }
func (ods SetDpiDeviceStructWithSetReportError) GetFeatureReport(b []byte) (int, error) {
	return 0, nil
}
func (ods SetDpiDeviceStructWithSetReportError) Read(b []byte) (int, error) { return 0, nil }
func (ods SetDpiDeviceStructWithSetReportError) SendFeatureReport(b []byte) (int, error) {
	return -1, nil
}
func (ods SetDpiDeviceStructWithSetReportError) Write(b []byte) (int, error) { return 0, nil }

type SetDpiDeviceStructOk struct{}

func (ods SetDpiDeviceStructOk) Close() error                           { return nil }
func (ods SetDpiDeviceStructOk) GetFeatureReport(b []byte) (int, error) { return 0, nil }
func (ods SetDpiDeviceStructOk) Read(b []byte) (int, error)             { return 0, nil }
func (ods SetDpiDeviceStructOk) SendFeatureReport(b []byte) (int, error) {
	return INTELLIMOUSE_PRO_SET_REPORT_LENGTH, nil
}
func (ods SetDpiDeviceStructOk) Write(b []byte) (int, error) { return 0, nil }

func Test_Driver_SetDpi(t *testing.T) {
	t.Parallel()
	t.Run("Test mouse and device nil error", func(t *testing.T) {
		d := Driver{}
		err := d.SetDpi(1000)
		assert.Error(t, err)
		assert.Equal(t, "Mouse Object or Device not instantiated", err.Error())
		d.mouse = &IntelliMousePro{}
		err = d.SetDpi(1000)
		assert.Equal(t, "Mouse Object or Device not instantiated", err.Error())
	})
	t.Run("Test Send feature report unexpected", func(t *testing.T) {
		d := Driver{}
		d.mouse = &IntelliMousePro{}
		d.device = SetDpiDeviceStructWithSetReportError{} // interface override
		err := d.SetDpi(1000)
		assert.Error(t, err)
	})
	t.Run("Test Send feature okay", func(t *testing.T) {
		d := Driver{}
		d.mouse = &IntelliMousePro{}
		d.device = SetDpiDeviceStructOk{} // interface override
		err := d.SetDpi(1000)
		assert.Nil(t, err)
	})
}

type ReadDpiSendFeatureReportWrongLength struct{}

func (ods ReadDpiSendFeatureReportWrongLength) Close() error                           { return nil }
func (ods ReadDpiSendFeatureReportWrongLength) GetFeatureReport(b []byte) (int, error) { return 0, nil }
func (ods ReadDpiSendFeatureReportWrongLength) Read(b []byte) (int, error)             { return 0, nil }
func (ods ReadDpiSendFeatureReportWrongLength) SendFeatureReport(b []byte) (int, error) {
	return -1, nil
}
func (ods ReadDpiSendFeatureReportWrongLength) Write(b []byte) (int, error) { return 0, nil }

type ReadDpiGetFeatureReportError struct{}

func (ods ReadDpiGetFeatureReportError) Close() error { return nil }
func (ods ReadDpiGetFeatureReportError) GetFeatureReport(b []byte) (int, error) {
	return 0, errors.New("Test")
}
func (ods ReadDpiGetFeatureReportError) Read(b []byte) (int, error) { return 0, nil }
func (ods ReadDpiGetFeatureReportError) SendFeatureReport(b []byte) (int, error) {
	return INTELLIMOUSE_PRO_SET_REPORT_LENGTH, nil
}
func (ods ReadDpiGetFeatureReportError) Write(b []byte) (int, error) { return 0, nil }

type ReadDpiGetFeatureReportOk struct{}

func (ods ReadDpiGetFeatureReportOk) Close() error { return nil }
func (ods ReadDpiGetFeatureReportOk) GetFeatureReport(b []byte) (int, error) {
	b[5] = 0x3e
	b[4] = 0x80
	return 0, nil
}
func (ods ReadDpiGetFeatureReportOk) Read(b []byte) (int, error) { return 0, nil }
func (ods ReadDpiGetFeatureReportOk) SendFeatureReport(b []byte) (int, error) {
	return INTELLIMOUSE_PRO_SET_REPORT_LENGTH, nil
}
func (ods ReadDpiGetFeatureReportOk) Write(b []byte) (int, error) { return 0, nil }

func Test_Driver_ReadDpi(t *testing.T) {
	t.Parallel()
	t.Run("Test mouse and device nil error", func(t *testing.T) {
		d := Driver{}
		err := d.SetDpi(1000)
		assert.Error(t, err)
		assert.Equal(t, "Mouse Object or Device not instantiated", err.Error())
		d.mouse = &IntelliMousePro{}
		err = d.SetDpi(1000)
		assert.Equal(t, "Mouse Object or Device not instantiated", err.Error())
	})
	t.Run("Test Feature report wrong len (trigger)", func(t *testing.T) {
		d := Driver{}
		d.mouse = &IntelliMousePro{}
		d.device = ReadDpiSendFeatureReportWrongLength{} // interface override
		_, err := d.ReadDpi()
		assert.Error(t, err)
		assert.Equal(t, "Error when reading report", err.Error())
	})

	t.Run("Test GetFeatureReport error", func(t *testing.T) {
		d := Driver{}
		d.mouse = &IntelliMousePro{}
		d.device = ReadDpiGetFeatureReportError{} // interface override
		_, err := d.ReadDpi()
		assert.Error(t, err)
		assert.Equal(t, "Error when reading report", err.Error())

	})
	t.Run("Test Read Dpi OK", func(t *testing.T) {
		d := Driver{}
		d.mouse = &IntelliMousePro{}
		d.device = ReadDpiGetFeatureReportOk{} // interface override
		rs, err := d.ReadDpi()
		assert.Nil(t, err)
		assert.Equal(t, uint16(16000), rs)

	})
}
