package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ValidateHex(t *testing.T) {
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
