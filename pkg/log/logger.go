package log

import (
	"io"
	"log"
	"os"
	"reflect"
)

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger
	writer  io.Writer
}

func NewLoggerOfObject(object any) *Logger {
	return NewLogger(reflect.TypeOf(object).String())
}

func NewLogger(name string) *Logger {
	writer := io.Writer(os.Stdout)
	flags := log.Ldate | log.Ltime

	return &Logger{
		debug:   log.New(writer, "DEBUG   ["+name+"]: ", flags),
		info:    log.New(writer, "INFO    ["+name+"]: ", flags),
		warning: log.New(writer, "WARNING ["+name+"]: ", flags),
		err:     log.New(writer, "ERROR   ["+name+"]: ", flags),
		writer:  writer,
	}
}

func (l *Logger) Debug(v ...any) {
	l.debug.Println(v...)
}

func (l *Logger) Info(v ...any) {
	l.info.Println(v...)
}

func (l *Logger) Warning(v ...any) {
	l.warning.Println(v...)
}

func (l *Logger) Error(v ...any) {
	l.err.Println(v...)
}

func (l *Logger) Fatal(v ...any) {
	l.err.Println(v...)
	os.Exit(1)
}

func (l *Logger) Debugf(format string, v ...any) {
	l.debug.Printf(format, v...)
}

func (l *Logger) Infof(format string, v ...any) {
	l.info.Printf(format, v...)
}

func (l *Logger) Warningf(format string, v ...any) {
	l.warning.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...any) {
	l.err.Printf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.err.Printf(format, v...)
	os.Exit(1)
}
