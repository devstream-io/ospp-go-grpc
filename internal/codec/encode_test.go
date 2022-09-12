package codec

import (
	"bytes"
	"github.com/devstream/ospp-go-grpc/shared"
	"testing"
)

func TestEncodeBytes(t *testing.T) {
	b := []byte{1, 2, 3}

	u, tp, err := Encode(b)
	if err != nil {
		t.Error(err)
	}

	if tp != shared.CodecTypeBytes {
		t.Error("type is not bytes")
	}

	if !bytes.Equal(b, u) {
		t.Error("bytes is not equal")
	}
}

func TestEncodeMap(t *testing.T) {
	m := map[string]interface{}{
		"key": "value",
	}

	u, tp, err := Encode(m)
	if err != nil {
		t.Error(err)
	}

	if tp != shared.CodecTypeMap {
		t.Error("type is not map")
	}

	if !bytes.Equal(u, []byte{0x81, 0xa3, 0x6b, 0x65, 0x79, 0xa5, 0x76, 0x61, 0x6c, 0x75, 0x65}) {
		t.Error("map is not equal")
	}
}

func TestEncodeUnsupportedType(t *testing.T) {
	_, _, err := Encode("test")
	if err == nil {
		t.Error("error should not be nil")
	}
}
