package codec

import (
	"fmt"
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"github.com/vmihailenco/msgpack/v5"
)

func Encode(v interface{}) ([]byte, pb.CodecType, error) {
	if v == nil {
		return nil, pb.CodecType_Bytes, nil
	}

	switch t := v.(type) {
	case map[string]interface{}:
		bytes, err := msgpack.Marshal(t)
		if err != nil {
			return nil, 0, err
		}
		return bytes, pb.CodecType_Map, nil
	case []byte:
		return t, pb.CodecType_Bytes, nil
	default:
		return nil, 0, fmt.Errorf("unsupported type: %v", t)
	}
}
