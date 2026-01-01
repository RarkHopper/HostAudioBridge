UNAME := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ifeq ($(UNAME),darwin)
  OSTYPE := darwin
else ifeq ($(UNAME),linux)
  OSTYPE := linux
else
  OSTYPE := windows
endif

# Build
.PHONY: build
build:
	go build -o bin/server ./cmd/server

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
	cp -n .env.$(OSTYPE).example .env.$(OSTYPE)
