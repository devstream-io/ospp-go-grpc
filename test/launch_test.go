package test

import (
	"context"
	"github.com/devstream/ospp-go-grpc/core"
	"github.com/devstream/ospp-go-grpc/plugin"
	"net"
	"testing"
	"time"
)

func TestCoreServeWithoutPlugin(t *testing.T) {
	c := core.New(token)
	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second)

	c.Shutdown()
}

func TestServeAndMountSuccessfully(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second), core.WithInterfaces(map[string][]string{
		"interface1": {"func1"},
	}))

	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Log(err)
		}
	}()

	p := plugin.New("plugin", "1.0.0", token)
	p.Handle("func1", func(_ plugin.Context) (interface{}, error) { return nil, nil })

	go func() {
		if err := p.Mount(context.Background(), "localhost", port); err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(time.Second)
	c.Shutdown()
}

func TestCoreServeWithListenedPort(t *testing.T) {
	lis, err := net.Listen("tcp", ":32100")
	if err != nil {
		t.Error(err)
	}
	defer func(lis net.Listener) {
		_ = lis.Close()
	}(lis)

	if err = core.New(token, core.WithPort(port)).Serve(context.Background()); err == nil {
		t.Error("port has benn listened, but no error")
	}
}

func TestMountWithMismatchedToken(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second))

	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Log(err)
		}
	}()

	p := plugin.New("plugin", "1.0.0", token+"1")

	if err := p.Mount(context.Background(), "localhost", port); err == nil {
		t.Errorf("mount with mismatched token, but no error")
	}

	c.Shutdown()
}

func TestMountWithSamePlugin(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second))

	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Log(err)
		}
	}()

	p := plugin.New("plugin", "1.0.0", token)

	go func() {
		if err := p.Mount(context.Background(), "localhost", port); err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(time.Second)

	if err := p.Mount(context.Background(), "localhost", port); err == nil {
		t.Errorf("mount with same plugin, but no error")
	}

	c.Shutdown()
}

func TestMountNotImplAnyInterface(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second), core.WithInterfaces(map[string][]string{
		"interface1": {"func1"},
	}))

	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Log(err)
		}
	}()

	p := plugin.New("plugin", "1.0.0", token)

	if err := p.Mount(context.Background(), "localhost", port); err == nil {
		t.Errorf("not impl any interface, but no error")
	}

	c.Shutdown()
}
