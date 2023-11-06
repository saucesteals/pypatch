package pypatch

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type GILState uintptr

type Python struct {
	h windows.Handle

	// int PyRun_SimpleString(const char *command) (https://docs.python.org/3/c-api/veryhigh.html#c.PyRun_SimpleString)
	procRunSimpleString uintptr

	// PyGILState_STATE PyGILState_Ensure() (https://docs.python.org/3/c-api/init.html#c.PyGILState_Ensure)
	procGILStateEnsure uintptr

	// void PyGILState_Release(PyGILState_STATE) (https://docs.python.org/3/c-api/init.html#c.PyGILState_Release)
	procGILStateRelease uintptr
}

func New(ctx context.Context) (*Python, error) {
	h, err := findPython(ctx)
	if err != nil {
		return nil, err
	}

	p := &Python{h: h}
	if err := p.init(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Python) Inject(code string) error {
	state, err := p.GILStateEnsure()
	if err != nil {
		return err
	}
	if err := p.RunSimpleString(code); err != nil {
		return err
	}

	return p.GILStateRelease(state)
}

func (p *Python) RunSimpleString(code string) error {
	codeBytes, err := windows.BytePtrFromString(code)
	if err != nil {
		return err
	}

	r1, _, e1 := syscall.SyscallN(p.procRunSimpleString, uintptr(unsafe.Pointer(codeBytes)))
	if e1 != 0 {
		return fmt.Errorf("could not run simple string: %w", err)
	}

	if r1 != 0 {
		return errors.New("python error executing code")
	}

	return nil
}

func (p *Python) GILStateEnsure() (GILState, error) {
	r1, _, e1 := syscall.SyscallN(p.procGILStateEnsure)
	if r1 == 0 {
		return 0, fmt.Errorf("could not ensure GIL state: %w", e1)
	}

	return GILState(r1), nil
}

func (p *Python) GILStateRelease(state GILState) error {
	_, _, e1 := syscall.SyscallN(p.procGILStateRelease, uintptr(state))
	if e1 != 0 {
		return fmt.Errorf("could not release GIL state: %w", e1)
	}
	return nil
}

func (p *Python) getProcAddress(name string) (uintptr, error) {
	return windows.GetProcAddress(p.h, name)
}

func (p *Python) init() error {
	var err error

	p.procRunSimpleString, err = p.getProcAddress("PyRun_SimpleString")
	if err != nil {
		return err
	}

	p.procGILStateEnsure, err = p.getProcAddress("PyGILState_Ensure")
	if err != nil {
		return err
	}

	p.procGILStateRelease, err = p.getProcAddress("PyGILState_Release")
	if err != nil {
		return err
	}

	return nil
}
