package core

import (
	"errors"
	"fmt"
	"github.com/devstream/ospp-go-grpc/internal/codec"
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"github.com/devstream/ospp-go-grpc/internal/util"
	"google.golang.org/protobuf/proto"
	"time"
)

type execResp struct {
	*pb.CommunicateExecResponse
}

func (c *Core) recvExecResp(resp *pb.CommunicateExecResponse) {
	exec, ok := c.execResp.Load(resp.ID)
	if !ok {
		return
	}
	ch := exec.(chan execResp)
	defer func() {
		if _, loaded := c.execResp.LoadAndDelete(resp.ID); loaded {
			close(ch)
		}
	}()

	r := execResp{
		CommunicateExecResponse: &pb.CommunicateExecResponse{
			ID:     resp.ID,
			Result: resp.Result,
			Type:   resp.Type,
			Err:    resp.Err,
		}}

	select {
	case ch <- r:
	default:
	}

}

// Call blocks until the func is executed or timeout
//
// args can be map[string]interface{} or []byte
func (c *Core) Call(plugin, version, funcName string, args interface{}) (Result, error) {
	p, ok := c.plugins.Load(util.GenKey(plugin, version))
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", plugin)
	}

	// reduce functions not found after transmission
	if !p.(*pluginInfo).funcs.Contains(funcName) {
		return nil, fmt.Errorf("func %s not found", funcName)
	}

	id := uint64(time.Now().UnixNano())

	// set result channel
	respCh := make(chan execResp, 0)
	c.execResp.Store(id, respCh)
	defer func() {
		if _, loaded := c.execResp.LoadAndDelete(id); loaded {
			close(respCh)
		}
	}()

	bytes, t, err := codec.Encode(args)
	if err != nil {
		return nil, err
	}
	b, err := proto.Marshal(&pb.CommunicateExecRequest{
		ID:       id,
		FuncName: funcName,
		Type:     t,
		Args:     bytes,
	})
	// failed to marshal
	if err != nil {
		return nil, err
	}
	if err = p.(*pluginInfo).comm.Send(&pb.CommunicateMsg{Type: pb.CommunicateType_ExecRequest, Data: b}); err != nil {
		return nil, err
	}

	// exec timeout
	timer := time.NewTimer(c.opts.execTimeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil, fmt.Errorf("exec %s.%s.%s timeout", plugin, version, funcName)
	case result := <-respCh:
		if result.Err != nil {
			return nil, errors.New(*result.Err)
		}

		union, err := codec.Decode(result.Result, result.Type)
		if err != nil {
			return nil, err
		}
		return &nativeResult{
			Union: union,
		}, nil
	}
}
