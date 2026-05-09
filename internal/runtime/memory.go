package runtime

import (
	"fmt"
	"unsafe"
)

type MemoryProbe struct{}

func (MemoryProbe) AddressOfString(s string) string {
	return fmt.Sprintf("0x%x", uintptr(unsafe.Pointer(unsafe.StringData(s))))
}

func (MemoryProbe) AddressOfBytes(b []byte) string {
	if len(b) == 0 {
		return "0x0"
	}
	return fmt.Sprintf("0x%x", uintptr(unsafe.Pointer(&b[0])))
}
