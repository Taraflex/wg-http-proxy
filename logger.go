package main

import (
	"log"
	"os"
)

var LogLevel = struct {
	FATAL uint64
	ERROR uint64
	WARN  uint64
	INFO  uint64
	DEBUG uint64
}{0b1, 0b11, 0b111, 0b1111, 0b11111}

type Logger struct {
	Level uint64
}

var stdout = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)
var stderr = log.New(os.Stderr, "WARN: ", log.LstdFlags)

func (l *Logger) Printf(format string, v ...interface{}) {
	if l.Level&LogLevel.DEBUG != 0 {
		stdout.Printf(format, v...)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.Level&LogLevel.WARN != 0 {
		stderr.Printf(format, v...)
	}
}
