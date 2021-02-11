package heaven

type LIST_ENTRY struct {
	Flink uint64
	Blink uint64
}

type PEB64 struct {
	Reserved          [24]byte
	LdrData           uint64
	ProcessParameters uint64
}

type PEB_LDR_DATA64 struct {
	Length                uint32
	Initialized           uint32
	SsHandle              uint64
	InLoadOrderModuleList LIST_ENTRY
}

type UNICODE_STRING_WOW64 struct {
	Length        uint16
	MaximumLength uint16
	Buffer        uint64
}

type UNICODE_STRING_WTF struct {
	Length        uint16
	MaximumLength uint16
	WtfIsThis     uint32
	Buffer        uint64
}

type ANSI_STRING_WOW64 struct {
	Length        uint16
	MaximumLength uint16
	WtfIsThis     uint32
	Buffer        uint64
}

type LDR_DATA_TABLE_ENTRY64 struct {
	InLoadOrderLinks           LIST_ENTRY
	InMemoryOrderLinks         LIST_ENTRY
	InInitializationOrderLinks LIST_ENTRY
	DllBase                    uint64
	EntryPoint                 uint64
	SizeOfImage                uint32
	Dummy                      uint64
	FullDllName                UNICODE_STRING_WOW64
	BaseDllName                UNICODE_STRING_WTF // [Length][Max][??extra 4 bytes??][Buffer]
}

type PROCESS_BASIC_INFORMATION64 struct {
	ExitStatus                   uint64
	PebBaseAddress               uint64
	AffinityMask                 uint64
	BasePriority                 uint64
	UniqueProcessId              uint64
	InheritedFromUniqueProcessId uint64
}
