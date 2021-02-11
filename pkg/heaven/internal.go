// +build 386

package heaven

import (
	"bytes"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/Binject/debug/pe"
	"golang.org/x/sys/windows"
)

var ntdllSize uint32

func getModuleHandle(module string) (uint64, error) {
	var h windows.Handle

	// Query 64-bit process information
	pInfo := PROCESS_BASIC_INFORMATION64{}
	size := uint32(unsafe.Sizeof(pInfo))

	h = GetSelfHandle()
	err := NtWow64QueryInformationProcess64(Handle(h), 0, windows.Pointer(unsafe.Pointer(&pInfo)), size, nil)
	if err != nil {
		return 0, fmt.Errorf("Could not get x64 module handle: NtWow64QueryInformationProcess64, %v", err)
	}
	windows.CloseHandle(h)

	// Read and build PEB64
	peb := PEB64{}

	h = GetSelfHandle()
	err = NtWow64ReadVirtualMemory64(Handle(h), pInfo.PebBaseAddress, windows.Pointer(unsafe.Pointer(&peb)), uint64(unsafe.Sizeof(peb)), nil)
	if err != nil {
		return 0, fmt.Errorf("Could not get x64 module handle: NtWow64ReadVirtualMemory64(peb), %v", err)
	}
	windows.CloseHandle(h)

	// Read and build ldr
	ldr := PEB_LDR_DATA64{}
	h = GetSelfHandle()
	err = NtWow64ReadVirtualMemory64(Handle(h), (peb.LdrData), windows.Pointer(unsafe.Pointer(&ldr)), uint64(unsafe.Sizeof(ldr)), nil)
	if err != nil {
		return 0, fmt.Errorf("Could not get x64 module handle: NtWow64ReadVirtualMemory64(head), %v", err)
	}
	windows.CloseHandle(h)

	// Read and build ldr data
	head := LDR_DATA_TABLE_ENTRY64{}
	head.InLoadOrderLinks.Flink = ldr.InLoadOrderModuleList.Flink

	lastEntry := peb.LdrData + 0x10
	for head.InLoadOrderLinks.Flink != lastEntry {

		h = GetSelfHandle()
		err = NtWow64ReadVirtualMemory64(Handle(h), head.InLoadOrderLinks.Flink, windows.Pointer(unsafe.Pointer(&head)), uint64(unsafe.Sizeof(head)), nil)
		if err != nil {
			return 0, fmt.Errorf("Could not get x64 module handle: NtWow64ReadVirtualMemory64(head loop), %v", err)
		}
		windows.CloseHandle(h)

		otherModLen := head.BaseDllName.Length / 2 //sizeof(wchar_t)
		if otherModLen != uint16(len(module)) {
			continue
		}

		// length match found, now check name
		// read BaseDllName.Buffer to buffer, convert from UTF16 to string

		var buff [256]uint16

		h = GetSelfHandle()
		err = NtWow64ReadVirtualMemory64(Handle(h), head.BaseDllName.Buffer, windows.Pointer(unsafe.Pointer(&buff)), uint64(head.BaseDllName.Length), nil)
		if err != nil {
			return 0, fmt.Errorf("Could not get x64 module handle: NtWow64ReadVirtualMemory64(otherModName), %v", err)
		}
		windows.CloseHandle(h)

		otherModName := windows.UTF16PtrToString(&buff[0])

		// check if module matches requested module, return or continue

		if otherModName == module {
			// hacky af
			if module == "ntdll.dll" {
				ntdllSize = head.SizeOfImage
			}
			return head.DllBase, nil
		}

	}

	return 0, fmt.Errorf("Could not get x64 module handle: module not found")
}

func getProcAddress(handle uint64, proc string) (uint64, error) {
	// get ldr proc base address
	ldrProcBaseAddress, err := getLdrProcedureAddress()
	if err != nil {
		return 0, err
	}

	// NtReadVirtualMemory
	// convert requested proc string to ANSI_STRING_WOW64
	buf := []byte(proc)
	pBuf := uint64((uintptr)(unsafe.Pointer(&buf[0])))
	aProc := ANSI_STRING_WOW64{
		Length:        uint16(len(buf)),
		MaximumLength: uint16(len(buf)),
		Buffer:        pBuf,
	}

	pProc := uint64((uintptr)(unsafe.Pointer(&aProc)))

	/*
		fmt.Println("callFunction(ldrProcBaseAddress, handle, pProc)")
		fmt.Printf("`- ldrProcBaseAddress:\t0x%016X\n", ldrProcBaseAddress)
		fmt.Printf("`- handle:\t\t0x%016X\n", handle)
		fmt.Printf("`- &proc:\t\t0x%016X\n", pProc)
		fmt.Printf("  `- proc.Length: 0x%04X\n", aProc.Length)
		fmt.Printf("  `- proc.Max:    0x%04X\n", aProc.MaximumLength)
		fmt.Printf("  `- proc.Buffer: 0x%016X\n", aProc.Buffer)
	*/

	// invoke LdrGetProcedureAddress passing proc unicode -> return Proc address
	var ret uint64
	fnret := callFunction(ldrProcBaseAddress,
		handle,
		pProc,
		0,
		uint64((uintptr)(unsafe.Pointer(&ret))),
	)

	if fnret > 0 {
		return 0, fmt.Errorf("Error in callFunction() ret: %v fnret: %v", ret, fnret)
	}

	return ret, nil
}

func getLdrProcedureAddress() (uint64, error) {
	var h windows.Handle
	var ntdll *pe.File

	ntdllBase, err := getModuleHandle("ntdll.dll")
	if err != nil {
		return 0, err
	}

	// read memory @ ntdllBase to size
	// TODO: clean this up
	var b [10000000]byte
	h = GetSelfHandle()
	err = NtWow64ReadVirtualMemory64(Handle(h), ntdllBase, windows.Pointer(unsafe.Pointer(&b)), uint64(ntdllSize), nil)
	if err != nil {
		return 0, fmt.Errorf("Could not get ldr procedure address: NtWow64ReadVirtualMemory64(ldr), %v", err)
	}
	windows.CloseHandle(h)

	buff := bytes.NewReader(b[:])

	ntdll, err = pe.NewFileFromMemory(buff)
	if err != nil {
		return 0, fmt.Errorf("Could not parse PE: %v", err)
	}

	exports, err := ntdll.Exports()
	if err != nil {
		return 0, err
	}

	for _, export := range exports {
		if export.Name == "LdrGetProcedureAddress" {
			return uint64(export.VirtualAddress) + ntdllBase, nil
		}
	}

	return 0, fmt.Errorf("LdrGetProcedureAddress not found")

}

func callFunction(proc uint64, args ...uint64) (errorcode uint64) {
	// probably should move this to a const and only allocate once instead of each time
	// via: https://github.com/JustasMasiulis/wow64pp
	shellcode := []byte{
		0x55,       // push ebp
		0x89, 0xE5, // mov ebp, esp

		0x83, 0xE4, 0xF0, // and esp, 0xFFFFFFF0

		// enter 64 bit mode
		0x6A, 0x33, // push 33h
		0xE8, 0x00, 0x00, 0x00, 0x00, // call
		0x83, 0x04, 0x24, 0x05, // add dword ptr [rsp],5
		0xCB, // retf

		0x67, 0x48, 0x8B, 0x4D, 16, // mov rcx, [ebp + 16] arg0
		0x67, 0x48, 0x8B, 0x55, 24, // mov rdx, [ebp + 24] arg1
		0x67, 0x4C, 0x8B, 0x45, 32, // mov r8,  [ebp + 32] arg2
		0x67, 0x4C, 0x8B, 0x4D, 40, // mov r9,  [ebp + 40] arg3

		0x67, 0x48, 0x8B, 0x45, 48, // mov rax, [ebp + 48] extSize

		0xA8, 0x01, // test al, 1
		0x75, 0x04, // jne 8, _no_adjust
		0x48, 0x83, 0xEC, 0x08, // sub rsp, 8
		// _no adjust:
		0x57,                         // push rdi
		0x67, 0x48, 0x8B, 0x7D, 0x38, // mov rdi, [ebp + 56] extArgs
		0x48, 0x85, 0xC0, // je _ls_e
		0x74, 0x16, 0x48, 0x8D, 0x7C, 0xC7, 0xF8, // lea rdi,[rdi+rax*8-8]
		// _ls:
		0x48, 0x85, 0xC0, // test rax, rax
		0x74, 0x0C, // je _ls_e
		0xFF, 0x37, // push [rdi]
		0x48, 0x83, 0xEF, 0x08, // sub rdi, 8
		0x48, 0x83, 0xE8, 0x01, // sub rax, 1
		0xEB, 0xEF, // jmp _ls
		// _ls_e:
		0x67, 0x8B, 0x7D, 0x40, // mov edi, [ebp + 64]
		0x48, 0x83, 0xEC, 0x20, // sub rsp, 0x20
		0x67, 0xFF, 0x55, 0x08, // call; [ebp + 0x8] func
		0x67, 0x89, 0x07, // mov [edi], eax
		0x67, 0x48, 0x8B, 0x4D, 0x30, // mov rcx, [ebp+48]
		0x48, 0x8D, 0x64, 0xCC, 0x20, // lea rsp,[rsp+rcx*8+0x20]
		0x5F, // pop rdi

		// exit 64 bit mode
		0xE8, 0, 0, 0, 0, 0xC7, 0x44, 0x24, 4, 0x23, 0, 0, 0, 0x83, 4, 0x24, 0xD, 0xCB,

		0x66, 0x8C, 0xD8, // mov ax, ds
		0x8E, 0xD0, // mov ss, eax

		0x89, 0xEC, // mov esp, ebp
		0x5D, // pop ebp
		0xC3, // ret
	}

	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")

	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_EXECUTE_READWRITE)

	if err != nil && err.Error() != "The operation completed successfully." {
		syscall.Exit(0)
	}

	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	RtlMoveMemory := ntdll.NewProc("RtlMoveMemory")
	_, _, err = RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))

	var ret uint32
	pret := (uint32)(uintptr(unsafe.Pointer(&ret)))

	// prep args array with minimum 4 elements
	l := 4
	if len(args) > 4 {
		l = len(args)
	}
	arrArgs := make([]uint64, l)
	for i, a := range args {
		arrArgs[i] = a
	}

	// how many additional args after args[3]
	extSize := uint64(l - 4)

	// if more args, set extArgs to pointer of 5th arg
	extArgs := uint64(0)
	if extSize > 0 {
		extArgs = (uint64)(uintptr(unsafe.Pointer(&arrArgs[4])))
	}

	/*
		fmt.Println("invoking heaven...")
		fmt.Printf("`- shellcode addr: 0x%08X\n", addr)
		fmt.Println("  `- args:")
		fmt.Printf("     `- proc:    0x%016X\n", proc)
		fmt.Printf("     `- arg0:    0x%016X\n", arrArgs[0])
		fmt.Printf("     `- arg1:    0x%016X\n", arrArgs[1])
		fmt.Printf("     `- arg2:    0x%016X\n", arrArgs[2])
		fmt.Printf("     `- arg3:    0x%016X\n", arrArgs[3])
		fmt.Printf("     `- extSize: 0x%016X\n", extSize)
		fmt.Printf("     `- extArgs: 0x%016X\n", extArgs)
		fmt.Printf("     `- &ret:    0x%08X\n", pret)
		fmt.Scanln()
	*/

	heaven(addr, proc, arrArgs[0], arrArgs[1], arrArgs[2], arrArgs[3], extSize, extArgs, pret)

	/*
		fmt.Println("we didn't crash <3")
		fmt.Printf("ret: 0x%016X\n", ret)
		fmt.Scanln()
	*/

	return uint64(ret)
}
