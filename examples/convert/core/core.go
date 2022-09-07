package main

import (
	"context"
	"fmt"
	"github.com/devstream/ospp-go-grpc/core"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"
)

var (
	port = 13001
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	c := core.New("123",
		core.WithLogLevel(core.LogLevelDebug),
		core.WithPort(port),
		core.WithInterfaces(map[string][]string{
			"Convert": {
				"EchoMap2Map", "EchoMap2Bytes", "EchoBytes2Map", "EchoBytes2Bytes",
			},
		}),
		core.WithExecReqChSize(5),
		core.WithExecTimeout(time.Second*5),
		core.WithServerOpts(grpc.WriteBufferSize(64*1024), grpc.ReadBufferSize(64*1024)),
		core.WithHealthTimeout(time.Second*15),
	)
	go func() {
		if err := c.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	cmd := exec.CommandContext(ctx, "../plugin/plugin", "--port", strconv.Itoa(port))
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// wait for plugin to start
	time.Sleep(time.Second * 3)

	call(c, "EchoMap2Map", map[string]interface{}{
		"Text": "hello",
	})
	call(c, "EchoMap2Bytes", map[string]interface{}{
		"Text": "hello",
	})
	call(c, "EchoBytes2Map", []byte("hello"))
	call(c, "EchoBytes2Bytes", []byte("hello"))

	<-ctx.Done()
}

func call(c *core.Core, name string, args interface{}) {
	start := time.Now()
	fmt.Printf("\ncall %s, args: %v\n", name, args)

	r, err := c.Call("convert", "v1", name, args)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("result: %v, err: %v, time: %v\n", r, err, time.Since(start))
}
