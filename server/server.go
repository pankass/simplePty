package main

import (
	"flag"
	"fmt"
	"github.com/creack/pty"
	"golang.org/x/term"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"unicode/utf8"
)

type TransText struct {
	r     io.Reader
	rawRw io.ReadWriter
}

func main() {

	var lport string
	flag.StringVar(&lport, "lport", "40000", "listen port")

	l, err := net.Listen("tcp", ":"+lport)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	fmt.Println("server listening on :" + lport)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("accepted new connection from " + conn.RemoteAddr().String())
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	ptmx, t, err := pty.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = ptmx.Close() }()

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go io.Copy(ptmx, os.Stdin)
	go io.Copy(os.Stdout, ptmx)

	transConn := NewTransConn(conn)

	go io.Copy(t, transConn)
	io.Copy(transConn, t)
}

func NewTransConn(conn net.Conn) *TransText {
	return &TransText{rawRw: conn}
}

// TransformText 编码转换
func (t *TransText) Read(buf []byte) (int, error) {
	if t.r != nil {
		return t.r.Read(buf)
	}
	n, err := t.rawRw.Read(buf)
	if err != nil {
		return n, err
	}
	// 适配windows下cmd的gbk编码
	if !utf8.Valid(buf[:n]) {
		buf, err = simplifiedchinese.GBK.NewDecoder().Bytes(buf)
		if err != nil {
			return n, err
		}
		t.r = transform.NewReader(t.rawRw, simplifiedchinese.GBK.NewDecoder())
	}
	return n, err
}

func (t *TransText) Write(buf []byte) (int, error) {
	return t.rawRw.Write(buf)
}
