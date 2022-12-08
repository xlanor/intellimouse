package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ValidateHex(t *testing.T) {
	t.Parallel()
	t.Run("Test invalid length", func(t *testing.T) {
		// use something we know will be valid char
		h := HexColor{}
		for i := 0; i < 8; i++ {
			if i != 3 && i != 6 {
				s := "#"
				for j := 0; j < i; j++ {
					s += "F"
				}
				h.Init(s)
				_, err := h.ValidateHex()
				assert.Error(t, err)
				assert.Equal(t, "Invalid hex string", err.Error())
			}
		}
	})
	t.Run("Test no hex code", func(t *testing.T) {
		h := HexColor{}
		teststr := "0FABCD"
		h.Init(teststr)
		_, err := h.ValidateHex()
		assert.Error(t, err)
		assert.Equal(t, "Invalid hex string", err.Error())
	})
	t.Run("Test Invalid Characters", func(t *testing.T) {
		h := HexColor{}
		teststr := "#G0FAB1"
		h.Init(teststr)
		_, err := h.ValidateHex()
		assert.Error(t, err)
		assert.Equal(t, "Invalid hex string", err.Error())
	})
	t.Run("Test three digit hex", func(t *testing.T) {
		h := HexColor{}
		teststr := "#ABC"
		h.Init(teststr)
		rs, err := h.ValidateHex()
		assert.Nil(t, err)
		assert.Equal(t, "AABBCC", rs)
	})
	t.Run("Test six digit hex", func(t *testing.T) {
		h := HexColor{}
		teststr := "#ABCDEF"
		h.Init(teststr)
		rs, err := h.ValidateHex()
		assert.Nil(t, err)
		assert.Equal(t, "ABCDEF", rs)
	})
}

func Test_BigEndianByteToHex(t *testing.T) {
	t.Parallel()
	t.Run("Test Big Endian Pad byte", func(t *testing.T) {
		input := []byte{255}
		ch := HexColor{}
		rs := ch.BigEndianByteToHex(input)
		assert.Equal(t, "#FF0000", rs)
	})
	t.Run("Test Big Endian Leading Zero", func(t *testing.T) {
		input := []byte{0, 0, 0}
		ch := HexColor{}
		rs := ch.BigEndianByteToHex(input)
		assert.Equal(t, "#000000", rs)
	})
}
