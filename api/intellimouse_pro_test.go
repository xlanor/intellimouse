package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	DPI = 128
)

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

func Test_SetDpi(t *testing.T) {
	t.Parallel()
	ims := IntelliMousePro{}
	t.Run("Test rounding logic", func(t *testing.T) {
		var testvar uint16
		testvar = 555
		rs := ims.SetDpiPayload(testvar)
		assert.Equal(t, uint16(600), GetUint16FromByte(rs, true))
	})
	t.Run("Test rounding logic exceed max", func(t *testing.T) {
		var testvar uint16
		testvar = 16999
		rs := ims.SetDpiPayload(testvar)
		assert.Equal(t, uint16(16000), GetUint16FromByte(rs, true))
	})
	t.Run("Test rounding logic exceed min", func(t *testing.T) {
		var testvar uint16
		testvar = 99
		rs := ims.SetDpiPayload(testvar)
		assert.Equal(t, uint16(200), GetUint16FromByte(rs, true))
	})
}

func Test_TriggerReadRequestPayload(t *testing.T) {
	t.Parallel()
	t.Run("Test trigger read request payload", func(t *testing.T) {
		ims := IntelliMousePro{}
		bs := ims.TriggerReadRequestPayload(INTELLIMOUSE_PRO_DPI_READ)
		// Validate byte slice
		assert.Equal(t, INTELLIMOUSE_PRO_SET_REPORT_LENGTH, len(bs))
		// Make sure all the properties are right.
		assert.Equal(t, uint8(INTELLIMOUSE_PRO_SET_REPORT), bs[0])
		assert.Equal(t, uint8(INTELLIMOUSE_PRO_DPI_READ), bs[1])
		assert.Equal(t, uint8(0x01), bs[2])
		// make sure rest of array is zeroed
		for i := 3; i < INTELLIMOUSE_PRO_SET_REPORT_LENGTH; i++ {
			assert.Equal(t, uint8(0x00), bs[i])
		}
	})
}

func Test_TriggerWriteDataRequestPayload(t *testing.T) {
	t.Parallel()
	t.Run("Test trigger write request payload", func(t *testing.T) {
		ims := IntelliMousePro{}

		data := []byte{255, 255, 255, 255}
		test_property := uint8(0x8F)
		bs := ims.TriggerWriteDataRequestPayload(test_property, data)

		// Validate byte slice
		assert.Equal(t, INTELLIMOUSE_PRO_SET_REPORT_LENGTH, len(bs))

		// Make sure all the properties are right.
		assert.Equal(t, uint8(INTELLIMOUSE_PRO_SET_REPORT), bs[0])
		assert.Equal(t, uint8(test_property), bs[1])
		assert.Equal(t, uint8(len(data)), bs[2])
		// make sure data is there,
		for i := 3; i < len(data); i++ {
			assert.Equal(t, uint8(0xFF), bs[i])
		}

		// make sure rest of array is zeroed
		for i := 3 + len(data); i < INTELLIMOUSE_PRO_SET_REPORT_LENGTH; i++ {
			assert.Equal(t, uint8(0x00), bs[i])
		}
	})
}

func Test_SetColorHexPayload(t *testing.T) {
	t.Parallel()
	t.Run("Test invalid color hex returns error", func(t *testing.T) {
		// most of it is already in Test_ValidateHex, just making sure error is propagated here
		invalid := "#GADF00"
		ims := IntelliMousePro{}
		_, err := ims.SetColorHexPayload(invalid)
		assert.Error(t, err)
	})

	t.Run("Test valid color hex returns byte array", func(t *testing.T) {
		// most of it is already in Test_ValidateHex, just making sure error is propagated here
		valid := "#FFFF00"
		ims := IntelliMousePro{}
		rs, err := ims.SetColorHexPayload(valid)
		assert.Nil(t, err)
		expected := []byte{255, 255, 0}
		assert.Equal(t, expected, rs)
		// also test leading zeroes
		rs, err = ims.SetColorHexPayload("#000000")
		expected = []byte{0, 0, 0}
		assert.Equal(t, expected, rs)
	})

}
