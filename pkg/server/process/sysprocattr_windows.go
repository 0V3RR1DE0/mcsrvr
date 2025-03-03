//go:build windows
// +build windows

package process

import "syscall"

// NewSysProcAttr returns a new SysProcAttr for Windows that hides the console window.
func NewSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x08000000, // 0x08000000 is CREATE_NO_WINDOW
		HideWindow:    true,
	}
}
