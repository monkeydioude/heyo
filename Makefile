.PHONY: all proto-gen

all: proto-gen

proto-gen:
	@protoc schampionne.proto --go_out=plugins=grpc:.