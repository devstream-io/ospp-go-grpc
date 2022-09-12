package codec

import (
	"bytes"
	"github.com/devstream/ospp-go-grpc/internal/pb"
	"github.com/devstream/ospp-go-grpc/shared"
	"testing"
)

func TestDecodeToBytes(t *testing.T) {
	b := []byte{1, 2, 3}

	u, err := Decode(b, pb.CodecType_Bytes)
	if err != nil {
		t.Error(err)
	}

	if u.Type() != shared.CodecTypeBytes {
		t.Error("type is not bytes")
	}

	if !bytes.Equal(b, u.Bytes()) {
		t.Error("bytes is not equal")
	}
}

func TestDecodeToMap(t *testing.T) {
	m := map[string]interface{}{
		"key": "value",
	}

	b, _, err := Encode(m)
	if err != nil {
		t.Error(err)
	}

	u, err := Decode(b, pb.CodecType_Map)
	if err != nil {
		t.Error(err)
	}

	if u.Type() != shared.CodecTypeMap {
		t.Error("type is not map")
	}

	if u.Map().GetString("key") != "value" {
		t.Error("map is not equal")
	}
}

func TestDecodeToUnsupportedType(t *testing.T) {
	_, err := Decode([]byte("test"), pb.CodecType_Map)
	if err == nil {
		t.Error("error should not be nil")
	}
}
