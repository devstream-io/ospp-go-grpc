package plugin

import (
	"context"
	"fmt"
	"github.com/devstream/ospp-go-grpc/internal/pb"
)

func (p *Plugin) recv(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := p.clients.comm.Recv()
			if err != nil {
				return fmt.Errorf("communicate recv error: %v", err)
			}

			switch msg.Type {
			case pb.CommunicateType_ExecRequest:
				go p.exec(msg.Data)
			}
		}
	}
}
