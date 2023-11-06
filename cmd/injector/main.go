package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/saucesteals/pypatch/inject"
	"golang.org/x/sys/windows"
)

var (
	program = "program.exe"
	dll     = "pypatch.dll"
)

func openProgram() (*exec.Cmd, windows.Handle, error) {
	programPath, err := filepath.Abs(program)
	if err != nil {
		return nil, 0, err
	}

	cmd := exec.Command(programPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return nil, 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	child, err := tryFindFirstChild(ctx, uint32(cmd.Process.Pid))
	cancel()
	if err != nil {
		return nil, 0, err
	}

	log.Println("Found child")

	h, err := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, child)
	if err != nil {
		return cmd, 0, err
	}

	return cmd, h, nil
}

func main() {
	cmd, handle, err := openProgram()
	if err != nil {
		panic(err)
	}

	log.Println("Opened")

	_, err = inject.Inject(handle, dll)
	if err != nil {
		panic(err)
	}

	log.Println("Injected")

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}
