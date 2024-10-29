//go:build windows
// +build windows

package ptyexec

import (
	"net"
	"os/exec"
)

// Exec windows下使用
func Exec(cmd *exec.Cmd, conn net.Conn) {
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}
