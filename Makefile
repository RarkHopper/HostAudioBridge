# Lint
.PHONY: check
check:
	go tool golangci-lint run ./...

# Format
.PHONY: format
format:
	gofmt -w .
	go tool goimports -w .
