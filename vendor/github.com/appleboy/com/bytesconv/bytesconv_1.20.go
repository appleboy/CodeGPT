//go:build go1.20
// +build go1.20

package bytesconv

import "unsafe"

// StringToBytes converts string to byte slice without a memory allocation.
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func StrToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts byte slice to string without a memory allocation.
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func BytesToStr(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
