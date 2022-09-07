package plugin

import (
	"fmt"
	"github.com/devstream/ospp-go-grpc/internal/codec"
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"google.golang.org/protobuf/proto"
)

func (p *Plugin) exec(data []byte) {
	req := pb.CommunicateExecRequest{}
	if err := proto.Unmarshal(data, &req); err != nil {
		p.Log.Errorf("communicate exec req unmarshal error: %v", err)
		return
	}

	defer p.recoverExec(&req)
	resp, err := p.execFunc(&req)

	if err != nil {
		t := err.Error()
		s := &t
		resp = &pb.CommunicateExecResponse{
			ID:     req.ID,
			Type:   req.Type,
			Result: nil,
			Err:    s,
		}
	}

	msg, err := proto.Marshal(resp)
	if err != nil {
		p.Log.Errorf("communicate marshal exec response error: %v", err)
		return
	}

	if err = p.clients.comm.Send(&pb.CommunicateMsg{
		Type: pb.CommunicateType_ExecResponse,
		Data: msg,
	}); err != nil {
		p.Log.Errorf("communicate send exec response error: %v", err)
		return
	}
}

func (p *Plugin) recoverExec(req *pb.CommunicateExecRequest) {
	if r := recover(); r != nil {
		p.Log.Errorf("exec func %s(%d) panic: %v", req.FuncName, req.ID, r)
		p.opts.onPanic(p, req.ID, req.FuncName, fmt.Errorf("%v", r))
		// TODO refactor
		t := fmt.Errorf("panic: %v", r).Error()
		msg, err := proto.Marshal(&pb.CommunicateExecResponse{
			ID:     req.ID,
			Type:   req.Type,
			Result: nil,
			Err:    &t,
		})
		if err != nil {
			p.Log.Errorf("communicate marshal exec response error: %v", err)
			return

		}

		if err = p.clients.comm.Send(&pb.CommunicateMsg{
			Type: pb.CommunicateType_ExecResponse,
			Data: msg,
		}); err != nil {
			p.Log.Errorf("communicate send exec response error: %v", err)
			return
		}
	}
}

func (p *Plugin) execFunc(req *pb.CommunicateExecRequest) (*pb.CommunicateExecResponse, error) {
	f, ok := p.handlers.Load(req.FuncName)
	if !ok {
		return nil, fmt.Errorf("func %s not found", req.FuncName)
	}

	union, err := codec.Decode(req.Args, req.Type)
	if err != nil {
		return nil, err
	}

	result, err := f.(HandlerFunc)(&nativeCtx{
		plugin: p,
		Union:  union,
	})
	if err != nil {
		return nil, err
	}

	bytes, t, err := codec.Encode(result)
	if err != nil {
		return nil, err
	}

	return &pb.CommunicateExecResponse{
		ID:     req.ID,
		Type:   t,
		Result: bytes,
		Err:    nil,
	}, nil
}
