//go:build windows

package ui

import (
	"os"
	"os/exec"
)

// restartSelf spawns a new copy of the binary and waits for it to exit so
// the terminal prompt does not return prematurely.  syscall.Exec is not
// available on Windows, so we use os/exec instead.
func restartSelf() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	cmd := exec.Command(exe, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return
	}
	_ = cmd.Wait()
}
