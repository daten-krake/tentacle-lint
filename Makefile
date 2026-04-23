.PHONY: build test vet fmt lint lint-check clean

BINARY=tentacle-lint
VERSION ?= dev

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/$(BINARY) ./cmd/tentacle-lint

test:
	go test -race ./... -v

vet:
	go vet ./...

fmt:
	gofmt -w .

lint: vet fmt

lint-check:
	@gofmt -l . | read -r; if [ $$? -eq 0 ]; then \
		echo "Files not formatted:"; \
		gofmt -l .; \
		exit 1; \
	fi
	go vet ./...

clean:
	rm -rf bin/