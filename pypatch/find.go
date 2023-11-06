package pypatch

import (
	"context"

	"golang.org/x/sys/windows"
)

var (
	// Current Active Python Versions (https://www.python.org/downloads/)
	pythonModules = []string{
		"python38.dll",
		"python39.dll",
		"python310.dll",
		"python311.dll",
		"python312.dll",
		"python313.dll",
	}
)

func findPython(ctx context.Context) (windows.Handle, error) {
	var moduleErr error
	i := 0

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			moduleName, err := windows.UTF16PtrFromString(pythonModules[i])
			if err != nil {
				return 0, err
			}

			var handle windows.Handle
			moduleErr = windows.GetModuleHandleEx(0, moduleName, &handle)
			if moduleErr != nil || handle == 0 {
				i = (i + 1) % len(pythonModules)
				continue
			}

			return handle, nil
		}
	}
}
