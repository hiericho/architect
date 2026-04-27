//go:build windows
// +build windows

package architect

import (
	"syscall"
)

// setTCPWindow sets the initial TCP receive window size for Windows.
func setTCPWindow(fd uintptr, window int) {
	_ = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, window)
}
