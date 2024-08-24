default: build test-run

build:
    goreleaser build --config ./.goreleaser-dev.yaml --single-target --snapshot --clean

test-run:
    ./dist/fabric render "document.hello" --source-dir ./examples/templates/basic_hello/ -v

format:
    go mod tidy
    ./gen_code.sh

format-extra: format
    gofumpt -w -extra .

lint: format
    golangci-lint run

test:
    go test -timeout 10s -race -short -v ./...

test-pretty:
    gotestsum --format dots-v2 -- -timeout 10s -race -short -v ./...

test-all:
    go test -timeout 5m -race -v ./...

test-e2e:
    go test -timeout 5m -race -v ./test/e2e/...

generate:
    ./codegen/gen_code.sh

generate-docs:
    ./codegen/gen_docs.sh
