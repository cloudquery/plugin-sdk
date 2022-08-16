.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: gen-proto
gen-proto:
	protoc --proto_path=. --go_out . --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/pb/base.proto internal/pb/source.proto internal/pb/destination.proto