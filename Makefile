.PHONY: build test lint clean install

build:
	go build -o bin/promtrace ./cmd/promtrace
test:
	go test ./... -race -count=1
lint:
	golangci-lint run ./...
clean:
	rm -rf bin/
install:
	go install ./cmd/promtrace
dev:
	go run ./cmd/promtrace wrap python test.py
