package main

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Logger struct {
	next       Interacter
	remoteAddr string
}

func NewLogger(next Interacter) Interacter {
	return &Logger{
		next: next,
	}
}

func (rec *Logger) Exec() {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"remoteAddr": rec.remoteAddr,
			"took":       time.Since(start),
		}).Info("connection established")
	}(time.Now())
	rec.next.Exec()
}
