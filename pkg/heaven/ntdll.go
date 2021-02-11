package heaven

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const ERROR_SUCCESS syscall.Errno = 0
const ptrSize = unsafe.Sizeof(uintptr(0))

type Handle = syscall.Handle

var (
	ntdll                            *windows.DLL
	ntWow64QueryInformationProcess64 *windows.Proc
	ntWow64ReadVirtualMemory64       *windows.Proc
)

func init() {
	var err error
	ntdll, err = windows.LoadDLL("ntdll.dll")
	if err == nil {
		ntWow64QueryInformationProcess64, _ = ntdll.FindProc("NtWow64QueryInformationProcess64")
		ntWow64ReadVirtualMemory64, _ = ntdll.FindProc("NtWow64ReadVirtualMemory64")
	}
}

func NtWow64QueryInformationProcess64(processHandle Handle, processInformationClass int32,
	processInformation windows.Pointer, processInformationLength uint32, returnLength *uint32) error {

	if ntWow64QueryInformationProcess64 == nil {
		return fmt.Errorf("ntWow64QueryInformationProcess64==nil")
	}

	r1, _, err := ntWow64QueryInformationProcess64.Call(uintptr(processHandle), uintptr(processInformationClass),
		uintptr(unsafe.Pointer(processInformation)), uintptr(processInformationLength),
		uintptr(unsafe.Pointer(returnLength)))

	if int(r1) < 0 {
		if err != ERROR_SUCCESS {
			return err
		} else {
			return syscall.EINVAL
		}
	}

	return nil
}

func NtWow64ReadVirtualMemory64(processHandle Handle, baseAddress uint64,
	bufferData windows.Pointer, bufferSize uint64, returnSize *uint64) error {

	if ntWow64ReadVirtualMemory64 == nil {
		return fmt.Errorf("ntWow64ReadVirtualMemory64==nil")
	}

	var r1 uintptr
	var err error

	// this shouldnt ever happen
	if ptrSize == 8 {
		r1, _, err = ntWow64ReadVirtualMemory64.Call(uintptr(processHandle), uintptr(baseAddress),
			uintptr(unsafe.Pointer(bufferData)), uintptr(bufferSize), uintptr(unsafe.Pointer(returnSize)))
	} else {
		r1, _, err = ntWow64ReadVirtualMemory64.Call(uintptr(processHandle),
			uintptr(baseAddress&0xFFFFFFFF),
			uintptr(baseAddress>>32),
			uintptr(unsafe.Pointer(bufferData)),
			uintptr(bufferSize),
			uintptr(0),
			uintptr(unsafe.Pointer(returnSize)))
	}

	if int(r1) < 0 {
		if err != ERROR_SUCCESS {
			return err
		} else {
			return syscall.EINVAL
		}
	}

	return nil
}
