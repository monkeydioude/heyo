.PHONY: all proto

all: proto

.PHONY: proto-go
proto-go: # Re-generate the .pb.go files from their .proto parent
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	protoc --go_out=. --go-grpc_out=. -I proto proto/heyo.proto

.PHONY: proto-rust
proto-rust:
	cd proto/rust && cargo build

.PHONY: proto
proto: proto-rust proto-go

.PHONY: dev
dev:
	cd cmd/heyo && go run .

.PHONY: test
test:
	@go test -v ./...