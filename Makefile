ifeq ($(OS),Windows_NT)
    protoc_install :=
else ifeq ($(shell uname -s),Darwin)
    protoc_install := brew install protobuf
else ifeq ($(shell uname -s),Linux)
    protoc_install := apt install -y protobuf-compiler
endif

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: gen-proto
gen-proto:
	$(protoc_install)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	protoc --proto_path=. --go_out . --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/pb/base.proto internal/pb/source.proto internal/pb/destination.proto