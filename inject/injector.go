package inject

import (
	"golang.org/x/sys/windows"
)

var (
	modKernel32            = windows.NewLazyDLL("kernel32.dll")
	procVirtualAllocEx     = modKernel32.NewProc("VirtualAllocEx")
	procCreateRemoteThread = modKernel32.NewProc("CreateRemoteThread")
	procLoadLibraryA       = modKernel32.NewProc("LoadLibraryA")
)

// Inject injects a DLL into a process.
func Inject(handle windows.Handle, dllPath string) (windows.Handle, error) {
	bytePtr, err := windows.BytePtrFromString(dllPath)
	if err != nil {
		return 0, err
	}

	injectedHandle, err := inject(uintptr(handle), bytePtr, uintptr(len(dllPath)+1))
	if err != nil {
		return 0, err
	}

	return windows.Handle(injectedHandle), nil
}

func inject(handle uintptr, dllPath *byte, lenDllPath uintptr) (uintptr, error) {
	alloc, _, err := procVirtualAllocEx.Call(handle, 0, lenDllPath, windows.MEM_RESERVE|windows.MEM_COMMIT, windows.PAGE_EXECUTE_READWRITE)
	if alloc == 0 {
		return 0, err
	}

	z := uintptr(0)
	err = windows.WriteProcessMemory(windows.Handle(handle), alloc, dllPath, lenDllPath, &z)
	if err != nil {
		return 0, err
	}

	injectedHandle, _, err := procCreateRemoteThread.Call(handle, 0, 0, procLoadLibraryA.Addr(), alloc, 0, 0)
	if injectedHandle == 0 {
		return 0, err
	}

	return injectedHandle, nil
}
