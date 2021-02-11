# Gopher Heaven

[<img src="img/gopher.png" width="150">](img/gopher.png)

__All gophers go to heaven__

gopherheaven is a Go implementation of the classic Heaven's Gate technique originally published by [roy g biv](https://vx-underground.org/archive/VxHeaven/lib/vrg02.html) on VX Heaven in 2009. gopherheaven can be used as an evasion technique to directly call 64-bit code from a 32-bit process.

[@C-Sto](https://github.com/C-Sto) already [went to Go hell](https://github.com/C-Sto/BananaPhone) 😈, but @aus went to heaven. 😇

## Usage

If you are familiar with GetModuleHandle, GetProcAddress, and Syscall on Windows, the process is largely the same. See [examples/](/examples) directory for more. The following example shows invoking 64-bit [NtReadVirtualMemory](http://undocumented.ntinternals.net/index.html?page=UserMode%2FUndocumented%20Functions%2FMemory%20Management%2FVirtual%20Memory%2FNtReadVirtualMemory.html)

```go
ntdll, err := heaven.GetModuleHandle("ntdll.dll")
if err != nil {
  log.Fatal(err)
}

fn, err := heaven.GetProcAddress(ntdll, "NtReadVirtualMemory")
if err != nil {
  log.Fatal(err)
}

h := (uint64)(heaven.GetSelfHandle())
i := 6
b := 3
var read uint64

errcode, err := heaven.Syscall(
  fn,
  h, 
  uint64(uintptr(unsafe.Pointer(&i))),
  uint64(uintptr(unsafe.Pointer(&b))),
  4,
  uint64(uintptr(unsafe.Pointer(&read)))
)
```

## Build

Make sure your architecture is set to `GOARCH=386` and that you are executing on x64 Windows system. gopherheaven does not currently support what I call reverse Heaven's Gate (executing 32-bit code from a 64-bit process).

## Background

There's already alot of great publications on Heaven's Gate, so I will just you defer to these resources:

- http://blog.rewolf.pl/blog/?p=102
- https://vx-underground.org/archive/VxHeaven/lib/vrg02.html
- http://www.alex-ionescu.com/?p=300
- https://www.malwaretech.com/2013/06/rise-of-dual-architecture-usermode.html

## Why

I asked myself several times. 

## Other References

- https://github.com/rwfpl/rewolf-wow64ext
- https://github.com/JustasMasiulis/wow64pp
- https://github.com/karinkasweet/Gopher-sticker-pack