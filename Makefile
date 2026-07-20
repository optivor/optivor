.PHONY: build test lint release test-e2e clean install uninstall

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

install: build
	@echo "Installing Optivor binary and systemd service..."
	sudo install -d /etc/optivor /usr/local/bin
	sudo install -m 0755 bin/optivor /usr/local/bin/optivor
	sudo install -m 0644 deploy/systemd/optivor.service /etc/systemd/system/optivor.service
	@if [ ! -f /etc/optivor/optivor.yaml ]; then sudo cp optivor.yaml.example /etc/optivor/optivor.yaml; fi
	sudo systemctl daemon-reload
	@echo "Optivor systemd unit installed successfully. Run 'sudo systemctl enable --now optivor' to start."

uninstall:
	@echo "Uninstalling Optivor service and binary..."
	-sudo systemctl stop optivor 2>/dev/null
	-sudo systemctl disable optivor 2>/dev/null
	sudo rm -f /etc/systemd/system/optivor.service
	sudo rm -f /usr/local/bin/optivor
	sudo systemctl daemon-reload
	@echo "Optivor uninstalled cleanly."

clean:
	rm -rf bin/ dist/ /tmp/optivor-cache
