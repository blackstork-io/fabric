default: build test-run

build:
    goreleaser build --config ./.goreleaser-dev.yaml --single-target --snapshot --clean

test-run:
    ./dist/fabric render "document.hello" --source-dir ./examples/templates/basic_hello/ -v

format:
    go mod tidy
    gofumpt -w .
    gci write --skip-generated -s standard -s default -s "prefix(github.com/blackstork-io/fabric)" .

format-extra: format
    gofumpt -w -extra .

lint: format
    golangci-lint run

test:
    go test -timeout 10s -race -short -v ./...

test-all:
    go test -timeout 5m -race -v ./...

test-e2e:
    go test -timeout 5m -race -v ./test/e2e/...

generate:
    go generate ./...

generate-docs:
    go run ./tools/docgen --version v0.0.0-dev --output ./docs/plugins/
