package main

import (
	"context"
	"github.com/devstream/ospp-go-grpc/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
)

var log *plugin.Logger

func main() {
	// kill signal means it's killed by the core with exec.Start()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill)
	defer cancel()

	p := plugin.New("my-plugin-name", "v1", "devstream",
		plugin.WithLogLevel(plugin.LogLevelInfo),
		plugin.WithDialOpts(grpc.WithTransportCredentials(insecure.NewCredentials())))

	p.Handle("Read", Read)
	p.Handle("Create", Create)
	p.Handle("Update", Update)
	p.Handle("Delete", Delete)

	go func() {
		_ = p.MountLocal(ctx)
	}()

	log = p.Log
	<-ctx.Done()
}

func Read(ctx plugin.Context) (interface{}, error) {
	return ctx.Map().Map(), nil
}

func Create(ctx plugin.Context) (interface{}, error) {
	return ctx.Map().Map(), nil
}

func Update(ctx plugin.Context) (interface{}, error) {
	return ctx.Map().Map(), nil
}

func Delete(ctx plugin.Context) (interface{}, error) {
	return ctx.Map().Map(), nil
}
