package nes

import (
	_ "fmt" // Remember to remove the blank identifier here.
	"strings"
)

func replaceAtIndex(str string, replacement byte, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
}

// ConvertToHex Converts a heximal number to a human-readable string
func ConvertToHex(n uint16, d uint8) string {
	var s string = strings.Repeat("0", int(d))
	for i := int(d) - 1; i >= 0; i, n = i-1, n>>4 {
		s = replaceAtIndex(s, "0123456789ABCDEF"[n&0xF], i)
	}
	return s
}
