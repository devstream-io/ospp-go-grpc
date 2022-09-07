package core

import "github.com/devstream/ospp-go-grpc/internal/pb"

type execReq struct {
	*pb.CommunicateExecRequest
}

type execResp struct {
	*pb.CommunicateExecResponse
}
