package plugin

import (
	"fmt"
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"google.golang.org/protobuf/proto"
	"log"
)

func (l *Logger) Debug(v ...interface{}) {
	l.log(pb.LogLevel_Debug, fmt.Sprint(v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.log(pb.LogLevel_Info, fmt.Sprint(v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.log(pb.LogLevel_Warn, fmt.Sprint(v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.log(pb.LogLevel_Error, fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.log(pb.LogLevel_Debug, fmt.Sprintf(format, v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.log(pb.LogLevel_Info, fmt.Sprintf(format, v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.log(pb.LogLevel_Warn, fmt.Sprintf(format, v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.log(pb.LogLevel_Error, fmt.Sprintf(format, v...))
}

func (l *Logger) log(level pb.LogLevel, message string) {
	if l.plugin.opts.logLevel > level {
		return
	}

	b, err := proto.Marshal(&pb.LogInfo{
		Type:    level,
		Message: message,
	})
	if err != nil {
		log.Println(err)
		return
	}

	err = l.plugin.clients.comm.Send(&pb.CommunicateMsg{
		Type: pb.CommunicateType_Log,
		Data: b,
	})
	if err != nil {
		// can't log to server, log to stdout
		log.Printf("can't send log to core: %v", err)
		return
	}
}
