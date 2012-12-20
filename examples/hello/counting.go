package main

import (
	"net"
	"sync"
	"log"
)

type CountingListener struct {
	net.Listener
	wg sync.WaitGroup
}

func NewCountingListener(l net.Listener) net.Listener {
	return &CountingListener{
		Listener: l,
	}
}

func (l *CountingListener) Accept() (c net.Conn, err error) {
	log.Print("Accept...")
	c, err = l.Listener.Accept()
	log.Print("Accepted: ", c, err)
	if err == nil {
		log.Print("+1")
		l.wg.Add(1)
	}
	c2 := &CountingConn{
		Conn: c,
		l: l,
	}
	return c2, err
}

func (l *CountingListener) Close() error {
	log.Print("Listener waiting...")
	l.wg.Wait()
	// TODO racy, have a bool in listener, prevent new Accepts
	log.Print("Listener wait done.")
	err := l.Listener.Close()
	return err
}

type CountingConn struct {
	net.Conn
	l *CountingListener
}

func (c *CountingConn) Close() error {
	err := c.Conn.Close()
	if err == nil {
		log.Print("-1")
		c.l.wg.Done()
	}
	return err
}
