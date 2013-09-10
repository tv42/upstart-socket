// listen-like-upstart is a standalone command that emulates the TCP
// file descriptor passing of Upstart. It is intended for testing
// software that will be deployed via Upstart in production.
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

var prog = filepath.Base(os.Args[0])

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s [IP]:PORT CMD [ARG..]\n", prog)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		Usage()
		os.Exit(2)
	}
	l, err := net.Listen("tcp", flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: cannot listen: %v\n", prog, err)
		os.Exit(1)
	}

	tcp := l.(*net.TCPListener)
	f, err := tcp.File()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: cannot get listening FD: %v\n", prog, err)
		os.Exit(1)
	}
	fd, err := syscall.Dup(int(f.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: cannot duplicate listening FD: %v\n", prog, err)
		os.Exit(1)
	}
	os.Setenv("UPSTART_FDS", strconv.FormatUint(uint64(fd), 10))
	exe, err := exec.LookPath(flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: cannot find executable: %v\n", prog, err)
		os.Exit(1)
	}
	args := flag.Args()[1:]
	args[0] = exe
	err = syscall.Exec(exe, args, os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: cannot execute command: %v\n", prog, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "%s: successful exec left us running!\n", prog)
	os.Exit(3)
}
