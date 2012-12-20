package upstart

import (
	"errors"
	"net"
	"os"
	"strconv"
)

func Listen() (net.Listener, error) {
	fd_s := os.Getenv("UPSTART_FDS")
	if fd_s == "" {
		return nil, errors.New("UPSTART_FDS not set in environment")
	}
	fd, err := strconv.ParseUint(fd_s, 10, 0)
	if err != nil {
		return nil, err
	}
	path := "/dev/fd/" + fd_s
	file := os.NewFile(uintptr(fd), path)
	l, err := net.FileListener(file)
	if err != nil {
		return nil, err
	}
	return l, nil
}
