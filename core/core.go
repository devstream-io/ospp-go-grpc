package core

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"sync"
)

type Core struct {
	impl     *impl
	token    string
	plugins  sync.Map      // map[string]*pluginInfo
	opts     options       // options
	server   *grpc.Server  // grpc server
	status   pb.CoreStatus // status
	cron     *cron.Cron    // health check
	execResp sync.Map      // exec response map : map[uint64]chan execResp
	cancel   func()        // global ctx cancel function
}

type impl struct {
	core *Core
}

type pluginInfo struct {
	name     string
	version  string
	health   int64
	shutdown chan struct{}
	impl     string // name of plugin implementation
	funcs    mapset.Set

	comm pb.Conn_CommunicateServer
}

type Interface map[string][]string

func New(token string, opts ...Option) *Core {
	// TODO execTimeout 应当小于等于 serverOpts 的超时时间
	c := Core{
		token:    token,
		plugins:  sync.Map{},
		status:   pb.CoreStatus_Stopped,
		execResp: sync.Map{},
		cron:     cron.New(),
		opts:     defaultOpts(),
	}

	for _, opt := range opts {
		opt.apply(&c.opts)
	}

	c.server = grpc.NewServer(c.opts.serverOpts...)

	i := impl{core: &c}
	c.impl = &i
	pb.RegisterConnServer(c.server, &i)
	return &c
}
