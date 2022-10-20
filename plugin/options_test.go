package plugin

import (
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestWithCallOpts(t *testing.T) {
	opts := defaultOpts()
	WithCallOpts(grpc.MaxCallRecvMsgSize(200*1024), grpc.MaxCallSendMsgSize(100*1024)).apply(&opts)
	if len(opts.callOpts) != 2 {
		t.Error("expect 2 call opts")
	}
}

func TestWithDialOpts(t *testing.T) {
	opts := defaultOpts()
	WithDialOpts(grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithReadBufferSize(100*1024)).apply(&opts)
	if len(opts.dialOpts) != 2 {
		t.Error("expect 2 dial opts")
	}
}

func TestWithHeartbeat(t *testing.T) {
	opts := defaultOpts()
	WithHeartbeat(time.Second).apply(&opts)
	if opts.heartbeat != time.Second {
		t.Error("expect 1s heartbeat")
	}
}

func TestWithLogLevel(t *testing.T) {
	opts := defaultOpts()
	WithLogLevel(LogLevelWarn).apply(&opts)
	if opts.logLevel != pb.LogLevel(LogLevelWarn) {
		t.Error("expect warn log level")
	}
}

func TestWithOnPanic(t *testing.T) {
	opts := defaultOpts()
	WithOnPanic(func(plugin *Plugin, execID uint64, funcName string, err error) {}).apply(&opts)
	if opts.onPanic == nil {
		t.Error("expect not nil on panic")
	}
}
