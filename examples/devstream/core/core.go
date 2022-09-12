package main

import (
	"context"
	"fmt"
	"github.com/devstream/ospp-go-grpc/core"
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

	c := core.New("devstream",
		core.WithLogLevel(core.LogLevelInfo),
		core.WithPort(port),
		core.WithInterfaces(map[string][]string{
			"Plugin": {
				"Read", "Create", "Update", "Delete",
			},
		}),
		core.WithExecTimeout(20*time.Minute),
	)
	go func() {
		_ = c.Serve(ctx)
	}()

	cmd := exec.CommandContext(ctx, "../plugin/plugin", "--port", strconv.Itoa(port))
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("waiting for interrupt/kill signal")
	<-ctx.Done()
}
