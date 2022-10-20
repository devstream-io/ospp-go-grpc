package core

import (
	"testing"
	"time"
)

func TestWithHealthTimeout(t *testing.T) {
	opts := defaultOpts()
	WithHealthTimeout(time.Second).apply(&opts)
	if opts.healthTimeout != time.Second {
		t.Error("expect 1s health timeout")
	}
}

func TestWithLogLevel(t *testing.T) {
	opts := defaultOpts()
	WithLogLevel(LogLevelWarn).apply(&opts)
	if opts.logLevel != LogLevelWarn {
		t.Error("expect warn log level")
	}
}

func TestWithExecReqChSize(t *testing.T) {
	opts := defaultOpts()
	WithExecReqChSize(5).apply(&opts)
	if opts.execReqChSize != 5 {
		t.Error("expect 5 exec req ch size")
	}
}

func TestWithExecTimeout(t *testing.T) {
	opts := defaultOpts()
	WithExecTimeout(time.Second * 10).apply(&opts)
	if opts.execTimeout != time.Second*10 {
		t.Error("expect 10s exec timeout")
	}
}

func TestWithPort(t *testing.T) {
	opts := defaultOpts()
	WithPort(32000).apply(&opts)
	if opts.port != 32000 {
		t.Error("expect 32000 port")
	}
}

func TestWithLogger(t *testing.T) {
	opts := defaultOpts()
	WithLogger(nil).apply(&opts)
	if opts.logger != nil {
		t.Error("expect nil logger")
	}
}

func TestWithInterfaces(t *testing.T) {
	interfaces := map[string][]string{
		"test":  {"test"},
		"test2": {"test2", "test3"},
	}
	opts := defaultOpts()
	WithInterfaces(interfaces).apply(&opts)
	if opts.interfaces["test"].Cardinality() != 1 || opts.interfaces["test2"].Cardinality() != 2 {
		t.Error("expect nil interfaces")
	}
}
