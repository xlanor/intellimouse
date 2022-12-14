package api

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Thanks to deankarn (go-playground/colors)
const (
	hexRegexString = "^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$"
)

var (
	hexRegex = regexp.MustCompile(hexRegexString)
)

type HexColor struct {
	hex_string string
}

func (h *HexColor) Init(hexString string) {
	h.hex_string = hexString
}

func (h *HexColor) ValidateHex() (string, error) {
	s := strings.ToUpper(h.hex_string)
	if !hexRegex.MatchString(s) {
		return "", errors.New("Invalid hex string")
	}
	// Strip the hash, we wont be sending that.
	s = s[1:]
	// if not, check the string for 3 digit hex strings
	// expand them to 6 digits. Per w3 spec (https://www.w3.org/TR/2001/WD-css3-color-20010305#colorunits)
	// Anything that is 3 digit, for example #FF0 is equivalent to #FFFF00 (double each character)
	if len(s) == 3 {
		var ret strings.Builder
		for _, v := range s {
			ret.WriteRune(v)
			ret.WriteRune(v)
		}
		return ret.String(), nil
	}
	return s, nil
}

func (h *HexColor) BigEndianByteToHex(input []byte) string {
	var rs strings.Builder
	rs.WriteString("#")
	for i := 0; i < len(input); i++ {
		// 02 to respect the 2 space
		rs.WriteString(fmt.Sprintf("%02X", input[i]))
	}
	if len(input) < 3 {
		// pad to 3
		for i := len(input); i < 3; i++ {
			rs.WriteString("00")
		}
	}
	return rs.String()

}
