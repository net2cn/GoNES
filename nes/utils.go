package nes

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"
)

// ConvertUint8ToString Converts a uint8 to a human-readable decimal string.
func ConvertUint8ToString(n uint8) string {
	return strconv.Itoa(int(n))
}

// ReplaceStringAtIndex Replace a char in string of given index and return a new string.
func ReplaceStringAtIndex(str string, replacement byte, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
}

// ConvertToHex Converts a heximal number to a human-readable string.
func ConvertToHex(n uint16, d uint8) string {
	var s string = strings.Repeat("0", int(d))
	for i := int(d) - 1; i >= 0; i, n = i-1, n>>4 {
		s = ReplaceStringAtIndex(s, "0123456789ABCDEF"[n&0xF], i)
	}
	return s
}

// IsPathExists Check if a path exists, returns true when exists.
func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	fmt.Printf("Exception has occurred in IsPahtExists(): %s", err)
	return false
}

// ConvertColorToUint32 Convert a color.RGBA struct to an uint32.
func ConvertColorToUint32(c color.RGBA) uint32 {
	bytes := []byte{c.R, c.G, c.B, c.A}
	return binary.LittleEndian.Uint32(bytes)
}
