.PHONY: all proto

all: proto

.PHONY: proto-go
proto-go: # Re-generate the .pb.go files from their .proto parent
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	protoc --go_out=. --go-grpc_out=. -I proto proto/heyo.proto

.PHONY: proto-rust
proto-rust:
	cd proto/rust && cargo build

.PHONY: proto-ts
proto-ts:
	protoc \
		--plugin=./node_modules/.bin/protoc-gen-ts_proto \
		-I ./proto \
		--ts_proto_out=./proto/ts \
		--ts_proto_opt=outputServices=grpc-js,env=node,esModuleInterop=true \
		./proto/heyo.proto


.PHONY: proto
proto: proto-rust proto-go

.PHONY: dev
dev:
	cd cmd/heyo && go run .

.PHONY: test
test:
	@go test -v ./...

DEFAULT_CERT_PATH=cmd/tools/test_server/certs
CERT_PATH ?= ${DEFAULT_CERT_PATH}

.PHONY: certs
certs:
	openssl req -newkey rsa:2048 -nodes -keyout ${CERT_PATH}/localhost.key -x509 -days 365 -out ${CERT_PATH}/localhost.crt -subj "/CN=localhost"