package plugin

import (
	"context"
	"flag"
	"fmt"
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"net"
	"strconv"
)

var _port = flag.Int("port", 0, "the port that core is listening on")

// MountLocal is a helper to mount to local core, the plugin should be launched by exec
//
// port should be passed by 'port' flag
func (p *Plugin) MountLocal(ctx context.Context) error {
	if !flag.Parsed() {
		flag.Parse()
	}
	return p.Mount(ctx, "localhost", *_port)
}

// Mount mounts to core
func (p *Plugin) Mount(ctx context.Context, host string, port int) error {
	defer func() {
		_ = p.Shutdown(UnmountExit, nil)
	}()

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	conn, err := grpc.DialContext(ctx, net.JoinHostPort(host, strconv.Itoa(port)), p.opts.dialOpts...)
	if err != nil {
		return err
	}

	p.conn = conn
	p.clients.conn = pb.NewConnClient(conn)

	commClient, err := p.clients.conn.Communicate(ctx, p.opts.callOpts...)
	if err != nil {
		return err
	}
	p.clients.comm = commClient

	// mount
	funcs := make([]string, 0)
	p.handlers.Range(func(key, value interface{}) bool {
		funcs = append(funcs, key.(string))
		return true
	})
	b, err := proto.Marshal(&pb.MountRequest{
		Token:     p.token,
		Name:      p.name,
		Version:   p.version,
		Functions: funcs,
	})
	if err != nil {
		return err
	}

	if err = p.clients.comm.Send(&pb.CommunicateMsg{Type: pb.CommunicateType_Mount, Data: b}); err != nil {
		return err
	}

	// start heartbeat
	_, err = p.cron.AddFunc(fmt.Sprintf("@every %ds", int(p.opts.heartbeat.Seconds())), p.heartbeat())
	if err != nil {
		return err
	}
	p.cron.Start()

	// set connected status
	p.status = pb.PluginStatus_Connected

	return p.recv(ctx)
}

// unmount fill in `msg` will print the reason for unmounting in core, if you do not need to pass nil
func (p *Plugin) unmount(reason UnmountReason, msg *string) error {
	b, err := proto.Marshal(&pb.UnmountRequest{
		Reason:  pb.UnmountReason(reason),
		Token:   p.token,
		Name:    p.name,
		Version: p.version,
		Msg:     msg,
	})
	if err != nil {
		return err
	}

	return p.clients.comm.Send(&pb.CommunicateMsg{Type: pb.CommunicateType_Unmount, Data: b})
}

// Shutdown unmount and closes the connection
func (p *Plugin) Shutdown(reason UnmountReason, msg *string) []error {
	p.cron.Stop()

	errs := make([]error, 0)

	if err := p.unmount(reason, msg); err != nil {
		errs = append(errs, err)
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	// set status to disconnected
	p.status = pb.PluginStatus_Disconnected

	p.cancel()

	return errs
}
