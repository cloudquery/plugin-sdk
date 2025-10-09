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

.PHONY: coverage
coverage:
	go test -timeout 15m -coverprofile=coverage.out.tmp ./...
	cat coverage.out.tmp | grep -vE "MockGen|codegen|mocks" > coverage.out
	rm coverage.out.tmp
	echo "| File | Function | Coverage |" > coverage.md
	echo "| --- | --- | --- |" >> coverage.md
	go tool cover -func=coverage.out | tail -n +2 | while read line; do \
		file=$$(echo $$line | awk '{print $$1}'); \
		func=$$(echo $$line | awk '{print $$2}'); \
		cov=$$(echo $$line | awk '{print $$3}'); \
		printf "| %s | %s | %s |\\n" "$$file" "$$func" "$$cov" >> coverage.md; \
	done
	rm coverage.out