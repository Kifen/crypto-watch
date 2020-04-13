OPTS?=GO111MODULE=on

build: ## Build binaries
	${OPTS} go build -o ./bin/crypto-watch ./cmd/exchange
	${OPTS} go build -o ./apps/binance ./cmd/apps/binance

clean: ## Clean compiled binaires
	rm -rf bin
	rm -rf apps

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/Kifen/crypto-watch ./pkg
	goimports -w -local github.com/Kifen/crypto-watch ./cmd

install-linters: ## Install linters
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.24.0
	golangci-lint --version

lint: ## Run linters. Use make install-linters first.


proto:
	protoc -I. --go_out=plugins=grpc:. \
	  pkg/proto/exchange.proto
