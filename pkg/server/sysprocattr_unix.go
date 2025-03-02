// +build !windows

package server

import "syscall"

func newSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true, // Start a new session to detach from the parent's process group
	}
}
