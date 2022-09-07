package plugin

import (
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"google.golang.org/grpc"
	"time"
)

type Option interface {
	apply(*options)
}

type options struct {
	dialOpts []grpc.DialOption
	callOpts []grpc.CallOption

	heartbeat time.Duration
	logLevel  pb.LogLevel
	onPanic   func(plugin *Plugin, execID uint64, funcName string, err error)
}

type option struct {
	f func(*options)
}

func (o *option) apply(do *options) {
	o.f(do)
}

func newOption(f func(options *options)) *option { return &option{f: f} }

func defaultOpts() options {
	return options{
		dialOpts:  make([]grpc.DialOption, 0),
		callOpts:  make([]grpc.CallOption, 0),
		heartbeat: time.Second * 10,
		logLevel:  pb.LogLevel_Info,
		onPanic:   func(plugin *Plugin, execID uint64, funcName string, err error) {},
	}
}

// WithHeartbeat set heartbeat, default is 10s
func WithHeartbeat(dur time.Duration) Option {
	return newOption(func(options *options) {
		options.heartbeat = dur
	})
}

// WithLogLevel set log level, default is Info
func WithLogLevel(level LogLevel) Option {
	return newOption(func(options *options) {
		options.logLevel = pb.LogLevel(level)
	})
}

func WithDialOpts(opts ...grpc.DialOption) Option {
	return newOption(func(options *options) {
		options.dialOpts = opts
	})
}

func WithCallOpts(opts ...grpc.CallOption) Option {
	return newOption(func(options *options) {
		options.callOpts = opts
	})
}

// WithOnPanic default is Log.Errorf
func WithOnPanic(f func(plugin *Plugin, execID uint64, funcName string, err error)) Option {
	return newOption(func(options *options) {
		options.onPanic = f
	})
}
