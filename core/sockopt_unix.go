//go:build !windows
// +build !windows

package architect

import (
	"syscall"
)

// setTCPWindow sets the initial TCP receive window size for Unix-like systems.
func setTCPWindow(fd uintptr, window int) {
	_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, window)
}
