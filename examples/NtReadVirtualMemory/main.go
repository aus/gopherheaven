package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/aus/gopherheaven/pkg/heaven"
)

func main() {
	ntdll, err := heaven.GetModuleHandle("ntdll.dll")
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		os.Exit(1)
	}

	fmt.Printf("Module 64-bit Handle: 0x%016x\n", ntdll)

	fn, err := heaven.GetProcAddress(ntdll, "NtReadVirtualMemory")
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		os.Exit(1)
	}

	fmt.Printf("NtReadVirtualMemory Addr: 0x%016x\n", ntdll)

	h := (uint64)(heaven.GetSelfHandle())
	i := 6
	b := 3
	var read uint64

	fmt.Printf("Setup\n")
	fmt.Printf("`- handle: 0x%016X\n", h)
	fmt.Printf("`- i: %v  @ 0x%016X\n", i, uint64(uintptr(unsafe.Pointer(&i))))
	fmt.Printf("`- b: %v  @ 0x%016X\n", b, uint64(uintptr(unsafe.Pointer(&b))))

	fmt.Printf("Invoking NtReadVirtualMemory in 64-bit mode to read from &i and write to &b\n")

	errcode, err := heaven.Syscall(fn, h, uint64(uintptr(unsafe.Pointer(&i))), uint64(uintptr(unsafe.Pointer(&b))), 4, uint64(uintptr(unsafe.Pointer(&read))))

	fmt.Printf("NtReadVirtualMemory result:\n")
	fmt.Printf("`- err: %v\n", err)
	fmt.Printf("`- errcode: %v\n", errcode)
	fmt.Printf("`- number of bytes read: %v\n", read)
	fmt.Printf("`- i: %v\n", i)
	fmt.Printf("`- b: %v\n", b)

}
