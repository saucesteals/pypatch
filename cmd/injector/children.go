package main

import (
	"context"
	"log"
	"syscall"
	"time"
	"unsafe"
)

func tryFindFirstChild(ctx context.Context, rootPid uint32) (uint32, error) {
	ticker := time.NewTicker(time.Millisecond * 150)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-ticker.C:
			child, err := findFirstChild(rootPid)
			if err != nil {
				log.Println("Waiting for child...")
				continue
			}

			return child, nil
		}
	}
}

func findFirstChild(parentPid uint32) (uint32, error) {
	snapshot, err := syscall.CreateToolhelp32Snapshot(uint32(syscall.TH32CS_SNAPPROCESS), 0)
	if err != nil {
		return 0, err
	}
	defer syscall.CloseHandle(snapshot)

	var processEntry syscall.ProcessEntry32
	processEntry.Size = uint32(unsafe.Sizeof(processEntry))

	err = syscall.Process32First(snapshot, &processEntry)
	if err != nil {
		return 0, err
	}

	for {
		err = syscall.Process32Next(snapshot, &processEntry)
		if err != nil {
			return 0, err
		}

		if processEntry.ParentProcessID == parentPid {
			break
		}
	}

	return processEntry.ProcessID, nil
}
