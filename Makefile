.PHONY: test
test:
	go test -tags=assert -race -count=100 ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: benchmark
benchmark:
	go test -bench=Benchmark -run="^$$" ./...

benchmark-ci:
	go install go.bobheadxi.dev/gobenchdata@v1.2.1
	go test -bench . -benchmem ./... -run="^$$" | gobenchdata --json bench.json
	rm -rf .delta.* && go run scripts/benchmark-delta/main.go bench.json