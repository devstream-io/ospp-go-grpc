package core

import (
	"github.com/devstream/ospp-go-grpc/internal/codec"
)

type Result interface {
	codec.Union
}

type nativeResult struct {
	codec.Union
}
