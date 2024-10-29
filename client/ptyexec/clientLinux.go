//go:build linux || ignore || darwin
// +build linux ignore darwin

package ptyexec

import (
	"github.com/creack/pty"
	"io"
	"net"
	"os/exec"
)

// Exec linux下使用
func Exec(cmd *exec.Cmd, conn net.Conn) {

	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.
	go io.Copy(conn, ptmx)
	go io.Copy(ptmx, conn)

	cmd.Wait()
}
