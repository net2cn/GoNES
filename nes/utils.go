package nes

import (
	"fmt"
	_ "fmt" // Remember to remove the blank identifier here.
	"os"
	"strconv"
	"strings"
)

// ConvertUint8ToString Converts a uint8 to a human-readable decimal string.
func ConvertUint8ToString(n uint8) string {
	return strconv.Itoa(int(n))
}

func replaceAtIndex(str string, replacement byte, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
}

// ConvertToHex Converts a heximal number to a human-readable string.
func ConvertToHex(n uint16, d uint8) string {
	var s string = strings.Repeat("0", int(d))
	for i := int(d) - 1; i >= 0; i, n = i-1, n>>4 {
		s = replaceAtIndex(s, "0123456789ABCDEF"[n&0xF], i)
	}
	return s
}

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
