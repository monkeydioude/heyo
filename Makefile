.PHONY: all proto-gen

all: proto-gen

proto-gen: # Re-generate the .pb.go files from their .proto parent
	@protoc --go_out=plugins=grpc:. *.proto
