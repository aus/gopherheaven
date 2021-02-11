package heaven

import (
	"math/rand"
	"testing"
	"unsafe"
)

func TestHeaven(t *testing.T) {
	ntdll, err := GetModuleHandle("ntdll.dll")
	if err != nil {
		t.Errorf("GetModuleHandle failed: %v", err)
	}

	fn, err := GetProcAddress(ntdll, "NtReadVirtualMemory")
	if err != nil {
		t.Errorf("GetProcAddress failed: %v", err)
	}

	h := (uint64)(GetSelfHandle())
	i := 6
	b := 20
	var read uint64

	for idx := 0; idx < 200; idx++ {
		read = 0
		i = rand.Intn(100)
		b = rand.Intn(100)
		for j := 0; j < 60; j++ {
			_, err = Syscall(fn, h, uint64(uintptr(unsafe.Pointer(&i))), uint64(uintptr(unsafe.Pointer(&b))), 4, uint64(uintptr(unsafe.Pointer(&read))))

			if read != 4 {
				t.Errorf("Syscall failed: NtReadVirtualMemory() expected to read 4 bytes but read %v", read)
				return
			}

			if i != b {
				t.Errorf("Syscall failed: NtReadVirtualMemory() buffer mismatch i=%v b=%v", i, b)
				return
			}

		}
	}

	if read != 4 {
		t.Errorf("Syscall failed: NtReadVirtualMemory() expected to read 4 bytes but read %v", read)
	}

	if i != b {
		t.Errorf("Syscall failed: NtReadVirtualMemory() buffer mismatch i=%v b=%v", i, b)
	}

}
