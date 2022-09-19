package tlog

import (
	"encoding/json"
	"fmt"
	"github.com/sarglase/tlog/hook"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"
)

type level int

const (
	TraceLevel level = iota
	InfoLevel
	DebugLevel
	ErrorLevel
	PanicLevel
	PrettyLevel
)

func (l level) WithPrefix() string {
	var prefix string
	switch l {
	case TraceLevel:
		prefix = "[trace]"
	case InfoLevel:
		prefix = "[info]"
	case DebugLevel:
		prefix = "[debug]"
	case ErrorLevel:
		prefix = "[error]"
	case PanicLevel:
		prefix = "[panic]"
	case PrettyLevel:
		prefix = "[pretty]"
	default:
		prefix = "[]"
	}
	return prefix

}

type option func(log *TLog)

func WithWriterOption(writer io.Writer) option {
	return func(log *TLog) {
		log.Writer = writer
	}
}

func WithHook(hook hook.Hook) option {
	return func(log *TLog) {
		log.Hook = hook
	}
}

type TLog struct {
	Name   string
	Writer io.Writer
	Hook   hook.Hook
	//logger
	Level level
}

func (tl *TLog) enable(l level) bool {

	if tl.Level < l {
		return false
	}
	return true
}

var defaultTLog = &TLog{
	Writer: os.Stdout,
	Hook:   nil,
	Level:  TraceLevel,
	Name:   " ",
}

func SetLevel(l level) {
	defaultTLog.Level = l
}

func SetName(name string) {
	defaultTLog.Name = " [" + name + "] "
}

func New(options ...option) {
	for _, o := range options {
		o(defaultTLog)
	}
}

func WithLevel(l level, v ...any) {
	if defaultTLog.Level > l {
		return
	}
	log(l, v...)
}

func WithLevelf(l level, v string, format ...interface{}) {
	if defaultTLog.Level > l {
		return
	}
	v = fmt.Sprintf(v, format...)
	log(l, v)
}

func Info(v ...any) {
	WithLevel(InfoLevel, v...)
}

func Debug(v ...any) {
	WithLevel(DebugLevel, v...)
}

func Error(v ...any) {
	WithLevel(ErrorLevel, v...)
}

func Panic(v ...any) {
	WithLevel(PanicLevel, v...)
}
func Pretty(v ...any) {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return
	}
	WithLevel(PrettyLevel, []any{string(b)}...)
}

func Infof(v string, format ...interface{}) {
	WithLevelf(InfoLevel, v, format...)
}

func Debugf(v string, format ...interface{}) {
	WithLevelf(DebugLevel, v, format...)
}

func Errorf(v string, format ...interface{}) {
	WithLevelf(ErrorLevel, v, format...)
}

func log(l level, v ...any) {
	_, file, line, _ := runtime.Caller(3)
	stackInfo := file + ":" + strconv.Itoa(line)
	logStr := fmt.Sprint(v...)
	now := time.Now().Format("2006-01-02 15:04:05")
	prefix := l.WithPrefix()
	var colorLevel color
	switch l {
	case TraceLevel:
		colorLevel = gray
	case InfoLevel:
		colorLevel = green
	case DebugLevel:
		colorLevel = yellow
	case ErrorLevel:
		colorLevel = red
	case PanicLevel:
		colorLevel = gray
	case PrettyLevel:
		colorLevel = gray
	default:
		colorLevel = blue
	}
	msg := fmt.Sprintf("%s \x1b[%dm%s%s%s %s \x1b[0m\n", now, colorLevel, prefix, defaultTLog.Name, stackInfo, logStr)
	defaultTLog.Writer.Write([]byte(msg))
	if defaultTLog.Hook != nil {
		defaultTLog.Hook.Write([]byte(msg))
	}

}
