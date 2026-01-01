# Build
.PHONY: build
build:
	go build -o bin/server ./cmd/server

# Run
.PHONY: up
up: build
	set -a && . ./.env && set +a && ./bin/server

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
