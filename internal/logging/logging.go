package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	default:
		return LevelInfo, fmt.Errorf("unknown log level: %s", s)
	}
}

type Logger struct {
	level  Level
	output io.Writer
	logger *log.Logger
}

func New(level string, format string) (*Logger, error) {
	lvl, err := ParseLevel(level)
	if err != nil {
		return nil, err
	}

	flags := log.Ldate | log.Ltime | log.Lmicroseconds
	if format == "text" {
		flags |= log.Lmsgprefix
	}

	return &Logger{
		level:  lvl,
		output: os.Stdout,
		logger: log.New(os.Stdout, "", flags),
	}, nil
}

func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
	l.logger.SetOutput(w)
}

func (l *Logger) SetLevel(level string) error {
	lvl, err := ParseLevel(level)
	if err != nil {
		return err
	}
	l.level = lvl
	return nil
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.log(LevelDebug, msg, args...)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.log(LevelInfo, msg, args...)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.log(LevelWarn, msg, args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= LevelError {
		l.log(LevelError, msg, args...)
	}
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(LevelError, msg, args...)
	os.Exit(1)
}

func (l *Logger) log(level Level, msg string, args ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", level.String())
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	l.logger.Printf("%s%s", prefix, msg)
}

var defaultLogger *Logger

func init() {
	var err error
	defaultLogger, err = New("info", "text")
	if err != nil {
		panic(err)
	}
}

func Default() *Logger {
	return defaultLogger
}

func SetLevel(level string) error {
	return defaultLogger.SetLevel(level)
}

func SetOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	defaultLogger.Fatal(msg, args...)
}
