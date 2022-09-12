package test

import (
	"context"
	"errors"
	"github.com/devstream/ospp-go-grpc/core"
	"github.com/devstream/ospp-go-grpc/plugin"
	"testing"
	"time"
)

func TestCoreExecTimeout(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second))

	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Log(err)
		}
	}()

	p := plugin.New("plugin", "1.0.0", token)
	p.Handle("func1", func(_ plugin.Context) (interface{}, error) { time.Sleep(2 * time.Second); return nil, nil })

	go func() {
		if err := p.Mount(context.Background(), "localhost", port); err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(time.Second)
	_, err := c.Call("plugin", "1.0.0", "func1", nil)
	if err == nil {
		t.Errorf("exec timeout, but no error")
	}

	c.Shutdown()
}

func TestCoreExecSuccessfully(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second))

	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Log(err)
		}
	}()

	p := plugin.New("plugin", "1.0.0", token)
	p.Handle("func1", func(_ plugin.Context) (interface{}, error) { time.Sleep(100 * time.Millisecond); return nil, nil })

	go func() {
		if err := p.Mount(context.Background(), "localhost", port); err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(time.Second)
	_, err := c.Call("plugin", "1.0.0", "func1", nil)
	if err != nil {
		t.Errorf("exec timeout, but no error")
	}

	c.Shutdown()
}

func TestCoreExecReturnError(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second))

	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Log(err)
		}
	}()

	p := plugin.New("plugin", "1.0.0", token)
	p.Handle("func1", func(_ plugin.Context) (interface{}, error) { return nil, errors.New("error") })

	go func() {
		if err := p.Mount(context.Background(), "localhost", port); err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(time.Second)
	_, err := c.Call("plugin", "1.0.0", "func1", nil)

	if err == nil {
		t.Errorf("expect error, got no error")
	}

	c.Shutdown()
}

func TestCoreExecReturnNoError(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second))

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
	_, err := c.Call("plugin", "1.0.0", "func1", nil)

	if err != nil {
		t.Errorf("expect no error, got error")
	}

	c.Shutdown()
}

func TestCoreExecWithBytesArgAndBytesResult(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second))

	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Log(err)
		}
	}()

	p := plugin.New("plugin", "1.0.0", token)
	p.Handle("func1", func(ctx plugin.Context) (interface{}, error) {
		return ctx.Bytes(), nil
	})

	go func() {
		if err := p.Mount(context.Background(), "localhost", port); err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(time.Second)
	result, err := c.Call("plugin", "1.0.0", "func1", []byte("hello"))

	if err != nil {
		t.Errorf("expect no error, got error")
	}

	if string(result.Bytes()) != "hello" {
		t.Errorf("expect bytes result, got no bytes result")
	}

	c.Shutdown()
}

func TestCoreExecWithBytesArgAndMapResult(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second))

	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Log(err)
		}
	}()

	p := plugin.New("plugin", "1.0.0", token)
	p.Handle("func1", func(ctx plugin.Context) (interface{}, error) {
		return map[string]interface{}{"hello": "world"}, nil
	})

	go func() {
		if err := p.Mount(context.Background(), "localhost", port); err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(time.Second)
	result, err := c.Call("plugin", "1.0.0", "func1", []byte("hello"))

	if err != nil {
		t.Errorf("expect no error, got error")
	}

	if result.Map().GetString("hello") != "world" {
		t.Errorf("expect map result, got no map result")
	}

	c.Shutdown()
}

func TestCoreExecWithMapArgAndBytesResult(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second))

	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Log(err)
		}
	}()

	p := plugin.New("plugin", "1.0.0", token)
	p.Handle("func1", func(ctx plugin.Context) (interface{}, error) {
		return []byte(ctx.Map().GetString("hello")), nil
	})

	go func() {
		if err := p.Mount(context.Background(), "localhost", port); err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(time.Second)
	result, err := c.Call("plugin", "1.0.0", "func1", map[string]interface{}{"hello": "world"})

	if err != nil {
		t.Errorf("expect no error, got error")
	}

	if string(result.Bytes()) != "world" {
		t.Errorf("expect bytes result, got no bytes result")
	}

	c.Shutdown()
}

func TestCoreExecWithMapArgAndMapResult(t *testing.T) {
	c := core.New(token, core.WithPort(port), core.WithExecTimeout(time.Second))

	go func() {
		if err := c.Serve(context.Background()); err != nil {
			t.Log(err)
		}
	}()

	p := plugin.New("plugin", "1.0.0", token)
	p.Handle("func1", func(ctx plugin.Context) (interface{}, error) {
		return map[string]interface{}{"hello": "world"}, nil
	})

	go func() {
		if err := p.Mount(context.Background(), "localhost", port); err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(time.Second)
	result, err := c.Call("plugin", "1.0.0", "func1", map[string]interface{}{"hello": "world"})

	if err != nil {
		t.Errorf("expect no error, got error")
	}

	if result.Map().GetString("hello") != "world" {
		t.Errorf("expect map result, got no map result")
	}

	c.Shutdown()
}
