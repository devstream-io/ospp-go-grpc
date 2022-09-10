package core

import (
	"fmt"
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"log"
)

type Logger interface {
	Log(prefix string, level LogLevel, v ...interface{})
	Logf(prefix string, level LogLevel, format string, v ...interface{})
}

type defaultLogger struct {
	level  LogLevel
	logger *log.Logger
}

var levelMap = map[LogLevel]string{
	LogLevelDebug: "DEBUG",
	LogLevelInfo:  "INFO",
	LogLevelWarn:  "WARN",
	LogLevelError: "ERROR",
}

func (l *defaultLogger) Log(prefix string, level LogLevel, v ...interface{}) {
	l.Logf(prefix, level, "%s", fmt.Sprint(v...))
}

func (l *defaultLogger) Logf(prefix string, level LogLevel, format string, v ...interface{}) {
	l.logger.Printf("%s [%s] %s", prefix, levelMap[level], fmt.Sprintf(format, v...))
}

func (i *impl) Log(server pb.Conn_LogServer) error {
	return nil
}
