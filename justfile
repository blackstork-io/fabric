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

generate:
    go generate ./...