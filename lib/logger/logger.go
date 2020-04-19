package logger

import (
	"log"
	"os"
)

const (
	defaultFlags = log.Ldate | log.Ltime
)

type Logger interface {
	Printf(string, ...interface{})
	Println(string)
	Errorf(string, ...interface{})
	Errorln(string)
}

type myLogger struct {
	stdoutLogger *log.Logger
	stderrLogger *log.Logger
}

func New() Logger {
	return &myLogger{
		stdoutLogger: log.New(os.Stdout, "", defaultFlags),
		stderrLogger: log.New(os.Stderr, "", defaultFlags),
	}
}

func (this *myLogger) Printf(format string, args ...interface{}) {
	this.stdoutLogger.Printf(format, args...)
}

func (this *myLogger) Println(message string) {
	this.stdoutLogger.Println(message)
}

func (this *myLogger) Errorf(format string, args ...interface{}) {
	this.stderrLogger.Printf(format, args...)
}

func (this *myLogger) Errorln(message string) {
	this.stderrLogger.Println(message)
}
