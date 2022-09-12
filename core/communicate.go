package core

import (
	"fmt"
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"google.golang.org/protobuf/proto"
	"strings"
	"time"
)

func (i *impl) Communicate(comm pb.Conn_CommunicateServer) error {
	bound := false
	plugin := &pluginInfo{}
	for {
		select {
		case <-plugin.shutdown:
			return fmt.Errorf("core force shutdown")
		default:
			recv, err := comm.Recv()
			if err != nil {
				return err
			}
			switch recv.Type {
			case pb.CommunicateType_Mount:
				if bound {
					continue
				}
				req := pb.MountRequest{}
				if err = proto.Unmarshal(recv.Data, &req); err != nil {
					return err
				}
				p, err := i.core.mount(&req, comm)
				if err != nil {
					return err
				}

				i.core.opts.logger.Logf("core", LogLevelInfo, "new plugin [%s.%s],impl [%s] interface,funcs: %v", p.name, p.version, p.impl, p.funcs.String())
				bound = true
				plugin = p
				plugin.health = time.Now().Unix() // init health time
			case pb.CommunicateType_Unmount:
				if !bound {
					continue
				}
				req := pb.UnmountRequest{}
				// parse failed, disconnect
				if err = proto.Unmarshal(recv.Data, &req); err != nil {
					return err
				}
				// unmount failed, disconnect
				i.core.opts.logger.Logf("core", LogLevelInfo, "unmount plugin %s.%s, %s:%v", req.Name, req.Version, pb.UnmountReason_name[int32(req.Reason)], req.Msg)
				return i.core.unmount(&req)
			case pb.CommunicateType_ExecResponse:
				if !bound {
					continue
				}
				resp := pb.CommunicateExecResponse{}
				if err = proto.Unmarshal(recv.Data, &resp); err != nil {
					return err
				}

				go i.core.recvExecResp(&resp)
			case pb.CommunicateType_Ping:
				if !bound {
					continue
				}
				plugin.health = time.Now().Unix()
			case pb.CommunicateType_Log:
				if !bound {
					continue
				}
				log := pb.LogInfo{}
				if err = proto.Unmarshal(recv.Data, &log); err != nil {
					return err
				}
				i.core.opts.logger.Log(strings.Join([]string{plugin.name, plugin.version}, "."), LogLevel(log.Type), log.Message)
			}
		}
	}
}
