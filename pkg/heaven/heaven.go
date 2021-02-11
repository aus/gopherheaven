// +build 386

package heaven

import (
	"fmt"

	"golang.org/x/sys/windows"
)

// GetModuleHandle returns a 64-bit handle to the specified module
func GetModuleHandle(module string) (uint64, error) {
	return getModuleHandle(module)
}

// GetProcAddress returns the 64-bit address of the exported function from the given 64-bit module handle
func GetProcAddress(handle uint64, proc string) (uint64, error) {
	return getProcAddress(handle, proc)
}

// Syscall initiates a 64-bit procedure at the specificed proc address
func Syscall(proc uint64, args ...uint64) (errcode uint32, err error) {
	errcode = (uint32)(callFunction(proc, args...))

	if errcode != 0 {
		err = fmt.Errorf("non-zero return from syscall")
	}
	return errcode, err
}

// GetSelfHandle returns a windows.Handle to the current process
func GetSelfHandle() windows.Handle {
	var h windows.Handle
	windows.DuplicateHandle(windows.CurrentProcess(), windows.CurrentProcess(), windows.CurrentProcess(), &h, 0, false, 0x00000002)
	return h
}
