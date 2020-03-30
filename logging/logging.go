package logging

import (
	"fmt"
	"log"
	"os"
)

type level int

const (
	TRACE level = iota
	DEBUG
	INFO
	WARN
	ERROR
)

func (l level) String() string {
	return [5]string{"[ TRACE ] ", "[ DEBUG ] ", "[ INFO  ] ", "[ WARN  ] ", "[ ERROR ] "}[l]
}

const flags = log.LstdFlags | log.Lmicroseconds | log.Lshortfile

type Logger struct {
	level       level
	traceLogger *log.Logger
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger(level level) *Logger {
	return &Logger{
		level:       level,
		traceLogger: log.New(os.Stdout, TRACE.String(), flags),
		debugLogger: log.New(os.Stdout, DEBUG.String(), flags),
		infoLogger:  log.New(os.Stdout, INFO.String(), flags),
		warnLogger:  log.New(os.Stdout, WARN.String(), flags),
		errorLogger: log.New(os.Stdout, ERROR.String(), flags),
	}
}

func NewLoggerFrom(specifier string) *Logger {
	return &Logger{
		level:       asLevel(specifier),
		traceLogger: log.New(os.Stdout, TRACE.String(), flags),
		debugLogger: log.New(os.Stdout, DEBUG.String(), flags),
		infoLogger:  log.New(os.Stdout, INFO.String(), flags),
		warnLogger:  log.New(os.Stdout, WARN.String(), flags),
		errorLogger: log.New(os.Stdout, ERROR.String(), flags),
	}
}

func (l *Logger) Trace(message string, values ...interface{}) {
	if l.level == TRACE {
		_ = l.traceLogger.Output(2, fmt.Sprintf(message, values...))
	}
}

func (l *Logger) Debug(message string, values ...interface{}) {
	if l.level <= DEBUG {
		_ = l.debugLogger.Output(2, fmt.Sprintf(message, values...))
	}
}

func (l *Logger) Info(message string, values ...interface{}) {
	if l.level <= INFO {
		_ = l.infoLogger.Output(2, fmt.Sprintf(message, values...))
	}
}

func (l *Logger) Warn(message string, values ...interface{}) {
	if l.level <= WARN {
		_ = l.warnLogger.Output(2, fmt.Sprintf(message, values...))
	}
}

func (l *Logger) Error(message string, values ...interface{}) {
	if l.level <= ERROR {
		_ = l.errorLogger.Output(2, fmt.Sprintf(message, values...))
	}
}

func (l *Logger) Fatal(message string, values ...interface{}) {
	_ = l.errorLogger.Output(2, fmt.Sprintf(message, values...))
	os.Exit(1)
}

func asLevel(specifier string) level {
	switch specifier {
	case "TRACE", "trace":
		return TRACE
	case "DEBUG", "debug":
		return DEBUG
	case "INFO", "info":
		return INFO
	case "WARNING", "WARN", "warning", "warn":
		return WARN
	case "ERROR", "error":
		return ERROR
	default:
		return INFO
	}
}
