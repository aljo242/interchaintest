default: help

.PHONY: help
help: ## Print this help message
	@echo "Available make commands:"; grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: interchaintest
interchaintest: gen ## Build interchaintest binary into ./bin
	go test -ldflags "-X github.com/strangelove-ventures/interchaintest/v5/internal/version.GitSha=$(shell git describe --always --dirty)" -c -o ./bin/interchaintest ./cmd/interchaintest

.PHONY: test
test: ## Run unit tests
	@go test -cover -short -race -timeout=60s ./...

.PHONY: docker-reset
docker-reset: ## Attempt to delete all running containers. Useful if interchaintest does not exit cleanly.
	@docker stop $(shell docker ps -q) &>/dev/null || true
	@docker rm --force $(shell docker ps -q) &>/dev/null || true

.PHONY: docker-mac-nuke
docker-mac-nuke: ## macOS only. Try docker-reset first. Kills and restarts Docker Desktop.
	killall -9 Docker && open /Applications/Docker.app

.PHONY: gen
gen: ## Run code generators
	go generate ./...

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --out-format=tab

lint-fix:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix --out-format=tab --issues-exit-code=0

.PHONY: lint lint-fix

format:
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.gw.go' | xargs go run mvdan.cc/gofumpt -w .
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.gw.go' | xargs go run github.com/client9/misspell/cmd/misspell -w
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.gw.go' | xargs go run golang.org/x/tools/cmd/goimports -w -local github.com/ingenuity-build/quicksilver
.PHONY: format