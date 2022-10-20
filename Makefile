.PHONY: protobuf
protobuf:
	protoc --version
	cd internal/pb && protoc -I . --go_out=plugins=grpc:. *.proto
