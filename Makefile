.PHONY: default ci unit integration lint vet gofmt staticcheck vulncheck


default:
	@echo 'Usage: make [unit|integration|lint|vet|gofmt|staticcheck|vulncheck|ci]'

# Fast, non-privileged checks suitable for any environment.
ci: vet gofmt staticcheck unit

# Compile-check + non-network unit tests.
unit:
	go test -run "^$$" ./...

# Privileged integration tests requiring a TUN/TAP kernel device.
# Run as: sudo make integration
integration:
	go test -v -timeout 120s ./...

vet:
	go vet ./...

gofmt:
	@out=$$(gofmt -s -l .); if [ -n "$$out" ]; then echo "$$out"; exit 1; fi

# staticcheck replaces the deprecated golint.
# Install once: go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck:
	staticcheck ./...

# Scan reachable vulnerabilities in dependencies.
# Install once: go install golang.org/x/vuln/cmd/govulncheck@latest
vulncheck:
	govulncheck ./...
