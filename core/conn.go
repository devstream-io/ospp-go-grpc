package core

import (
	"context"
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"github.com/devstream/ospp-go-grpc/internal/util"
	"net"
)

func (c *Core) Serve(ctx context.Context) error {
	defer func() {
		// if errors occur, close the server
		if c.status == pb.CoreStatus_Launched {
			return
		}
		c.Shutdown()
	}()

	lc := net.ListenConfig{}
	lis, err := lc.Listen(ctx, "tcp", fmt.Sprintf(":%d", c.opts.port))
	if err != nil {
		return err
	}

	// TODO attention
	// health check,only internal is the half of the healthTimeout
	_, err = c.cron.AddFunc(fmt.Sprintf("@every %ds", int(c.opts.execTimeout.Seconds()/2)), c.healthCheck())
	if err != nil {
		return err
	}
	c.cron.Start()
	defer c.cron.Stop()

	c.status = pb.CoreStatus_Launched

	return c.server.Serve(lis)
}

func (c *Core) ShutdownPlugin(plugin, version string) error {
	if c.status != pb.CoreStatus_Launched {
		return fmt.Errorf("core is not launched")
	}

	key := util.GenKey(plugin, version)

	p, ok := c.plugins.Load(key)
	if !ok {
		return fmt.Errorf("plugin %s not found", key)
	}

	close(p.(*pluginInfo).shutdown)
	c.plugins.Delete(key)

	return nil
}

// Shutdown core
func (c *Core) Shutdown() {
	c.cancel()
	c.cron.Stop()

	c.server.GracefulStop()

	c.status = pb.CoreStatus_Stopped
}

func (c *Core) bind(req *pb.BindRequest, comm pb.Conn_CommunicateServer) (*pluginInfo, error) {
	// invalid token, disconnect
	if req.Token != c.token {
		return nil, errors.New("invalid token")
	}

	// must impl only one of the interfaces
	funcs := mapset.NewSet()
	for _, f := range req.Functions {
		funcs.Add(f)
	}
	implName := ""
	if c.opts.interfaces != nil {
		impls := 0
		for name, set := range c.opts.interfaces {
			if funcs.IsSuperset(set) {
				impls++
				implName = name
			}
		}
		if impls != 1 {
			return nil, fmt.Errorf("must implement only one of the interfaces")
		}
	}

	key := util.GenKey(req.Name, req.Version)
	if _, ok := c.plugins.Load(key); ok {
		// 已存在插件断开连接
		return nil, fmt.Errorf("plugin %s.%s is exists", req.Name, req.Version)
	}

	info := pluginInfo{
		name:     req.Name,
		version:  req.Version,
		health:   0,
		shutdown: make(chan struct{}, 0),
		comm:     comm,
		impl:     implName,
		funcs:    funcs,
	}
	c.plugins.Store(key, &info)
	return &info, nil
}

func (c *Core) unbind(req *pb.UnbindRequest) error {
	if c.token != req.Token {
		return errors.New("invalid token")
	}
	key := util.GenKey(req.Name, req.Version)
	if _, ok := c.plugins.Load(key); !ok {
		return fmt.Errorf("plugin %s.%s is not exists", req.Name, req.Version)
	}
	c.plugins.Delete(key)
	return nil
}
