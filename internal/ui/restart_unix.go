//go:build !windows

package ui

import (
	"os"
	"syscall"
)

// restartSelf replaces the current process image with a fresh copy of the
// binary using syscall.Exec.  Because exec(2) overwrites the process in-place,
// the terminal sees the same PID and Ctrl-C continues to work normally.
// On failure the caller is expected to call fyneApp.Quit() so the user can
// restart manually.
func restartSelf() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	_ = syscall.Exec(exe, os.Args, os.Environ())
}
