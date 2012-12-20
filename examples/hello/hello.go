package main

import (
	"github.com/tv42/upstart-socket"
	"log"
	"net/http"
	"strings"
	"time"
//	"sync"
)

func hello(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Hello, world!\n"))
}

func main() {
	l, err := upstart.Listen()
	if err != nil {
		log.Fatalf("Cannot listen on Upstart-provided socket: %s", err)
	}
	// wg := sync.WaitGroup{}
	l2 := NewCountingListener(l)
	s := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
	// 	wg.Add(1)
	// 	defer wg.Done()
	// 	hello(rw, req)
	// })
	http.HandleFunc("/", hello)

	go func() {
		time.Sleep(time.Second / 10)
		log.Print("Deciding to stop.")
		l2.Close()
		log.Print("Listener closed.")
	}()

// xtg: smw_: keepalives are "fine" if it's just graceful termination, if you don't care how long that takes (particularly if you're just restarting the go http server), since it's just a matter of closing/freeing up the socket for new listens -- but the existing client conns can remain open
// smw_: Tv_, ah! you can write a transport to kill keepalives!
// smw_: Tv_, nm, you can implement it without forking the entire lib :-)
// smw_: xtg, ok

	log.Print("Serving...")
	err = s.Serve(l2)
	if err != nil {
		// kludge.. https://groups.google.com/forum/?fromgroups#!topic/golang-nuts/RMQNS2KA9CY%5B1-25%5D
		msg := err.Error()
		if msg == "use of closed network connection" ||
			(strings.HasPrefix(msg, "accept tcp ") &&
			strings.HasSuffix(msg, ": use of closed network connection")) {
			log.Print("Waiting...")
//			wg.Wait()
			log.Print("Stopping...")
		} else {
			log.Fatalf("Serving failed: %s\n", err)
		}
	}
}
