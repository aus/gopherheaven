// +build !386

package heaven

import (
	"fmt"

	"golang.org/x/sys/windows"
)

// GetModuleHandle returns a 64-bit handle to the specified module
func GetModuleHandle(module string) (uint64, error) {
	return 0, fmt.Errorf("unimplemented for this os/arch")
}

// GetProcAddress returns the 64-bit address of the exported function from the given 64-bit module handle
func GetProcAddress(handle uint64, proc string) (uint64, error) {
	return 0, fmt.Errorf("unimplemented for this os/arch")
}

// Syscall initiates a 64-bit procedure at the specificed proc address
func Syscall(proc uint64, args ...uint64) (errcode uint32, err error) {
	return 0, fmt.Errorf("unimplemented for this os/arch")
}

// GetSelfHandle returns a windows.Handle to the current process
func GetSelfHandle() windows.Handle {
	return 0
}
