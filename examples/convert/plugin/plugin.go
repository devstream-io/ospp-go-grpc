package main

import (
	"fmt"
	"github.com/devstream/ospp-go-grpc/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	p := plugin.New("convert", "v1", "123",
		plugin.WithLogLevel(plugin.LogLevelDebug),
		plugin.WithHeartbeat(10*time.Second),
		plugin.WithDialOpts(grpc.WithTransportCredentials(insecure.NewCredentials())))

	p.Handle("EchoMap2Map", EchoMap2Map)
	p.Handle("EchoMap2Bytes", EchoMap2Bytes)
	p.Handle("EchoBytes2Map", EchoBytes2Map)
	p.Handle("EchoBytes2Bytes", EchoBytes2Bytes)
	p.Handle("Panic", Panic)

	if err := p.MountLocal(); err != nil {
		log.Println(err)
		return
	}
}

func EchoMap2Map(ctx plugin.Context) (interface{}, error) {
	text := ctx.Map().GetString("Text")
	ctx.L().Debugf("echo|arg:map|result:map|arg:%v", text)
	return map[string]interface{}{
		"Text": text,
	}, nil
}

func EchoMap2Bytes(ctx plugin.Context) (interface{}, error) {
	text := ctx.Map().GetString("Text")
	ctx.L().Debugf("echo|arg:map|result:bytes|arg:%v", text)
	return []byte(text), nil
}

func EchoBytes2Map(ctx plugin.Context) (interface{}, error) {
	text := ctx.Bytes()
	ctx.L().Debugf("echo|arg:bytes|result:map|arg:%v", text)
	return map[string]interface{}{
		"Text": string(text),
	}, nil
}

func EchoBytes2Bytes(ctx plugin.Context) (interface{}, error) {
	text := ctx.Bytes()
	ctx.L().Debugf("echo|arg:bytes|result:bytes|arg:%v", text)
	return text, nil
}

func Panic(ctx plugin.Context) (interface{}, error) {
	ctx.L().Debug("I will panic")
	panic(fmt.Errorf("panic info"))
}
