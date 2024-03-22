.PHONY: test
test:
	go test -tags=assert -race ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: benchmark
benchmark:
	go test -bench=Benchmark -run="^$$" ./... | grep -v 'BenchmarkWriterMemory/'
	go test -bench=BenchmarkWriterMemory -run="^$$" ./writers/

benchmark-ci:
	go install go.bobheadxi.dev/gobenchdata@v1.2.1
	go test -bench . -benchmem ./... -run="^$$" | grep -v 'BenchmarkWriterMemory/' | gobenchdata --json bench.json
	rm -rf .delta.* && go run scripts/benchmark-delta/main.go bench.json
