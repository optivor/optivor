.PHONY: build test lint release test-e2e clean

build:
	go build -o bin/optivor ./cmd/optivor

test:
	go test ./... -race -cover

lint:
	go vet ./...
	@if command -v golangci-lint > /dev/null; then golangci-lint run; else echo "golangci-lint not installed, skipping"; fi

test-e2e:
	@echo "Running E2E tests..."
	go test -v ./test/e2e/...

release:
	goreleaser release --snapshot --clean

clean:
	rm -rf bin/ dist/ /tmp/optivor-cache
