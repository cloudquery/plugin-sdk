.PHONY: test
test:
	go test -race ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: gen-proto
gen-proto:
	protoc --proto_path=. --go_out . --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/pb/base.proto internal/pb/source.proto internal/pb/destination.proto

.PHONY: gen-proto
benchmark:
	go test -bench=Benchmark -run="^$$" ./...

benchmark-ci:
	go install go.bobheadxi.dev/gobenchdata@v1.2.1
	go test -bench . -benchmem ./... -run="^$$" | gobenchdata --json bench.json
	rm -rf .delta.* && go run scripts/benchmark-delta/main.go bench.json