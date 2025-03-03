//go:build !windows
// +build !windows

package process

import "syscall"

// NewSysProcAttr returns a new SysProcAttr for Unix-like systems
func NewSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		// Setsid creates a new session and detaches from the controlling terminal
		Setsid: true,
		// Pgid of 0 means a new process group is created
		Pgid: 0,
		// Foreground false indicates it should run in the background
		Foreground: false,
	}
}
