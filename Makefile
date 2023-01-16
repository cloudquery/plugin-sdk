.PHONY: test
test:
	go test -race ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: gen-proto
gen-proto:
	protoc --proto_path=. --go_out . --go_opt=module="github.com/cloudquery/plugin-sdk" --go-grpc_out=. --go-grpc_opt=module="github.com/cloudquery/plugin-sdk" cloudquery/base/v0/base.proto cloudquery/destination/v0/destination.proto cloudquery/source/v0/source.proto
	protoc --proto_path=. --go_out . --go_opt=module="github.com/cloudquery/plugin-sdk" --go-grpc_out=. --go-grpc_opt=module="github.com/cloudquery/plugin-sdk" cloudquery/source/v1/source.proto

.PHONY: benchmark
benchmark:
	go test -bench=Benchmark -run="^$$" ./...

benchmark-ci:
	go install go.bobheadxi.dev/gobenchdata@v1.2.1
	go test -bench . -benchmem ./... -run="^$$" | gobenchdata --json bench.json
	rm -rf .delta.* && go run scripts/benchmark-delta/main.go bench.json