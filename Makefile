# Build Server
.PHONY: build-server
build-server:
	go build -o bin/server ./cmd/server

# Build CLI
.PHONY: build-cli
build-cli:
	go build -o bin/hab-cli ./cmd/cli

.PHONY: build
build: build-server build-cli

# Run Server
.PHONY: up
up: build
	./bin/server

# Run CLI
.PHONY: cli
cli: build-cli
	./bin/hab-cli

# Test
.PHONY: test
test:
	go test ./... -v

# Lint
.PHONY: check
check:
	go tool golangci-lint run ./...

# Format
.PHONY: format
format:
	gofmt -w .
	go tool goimports -w .

# 環境変数ファイル初期化
.PHONY: env-init
env-init:
	cp -n .env.example .env
