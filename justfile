default: build build-plugins test-run

build:
    go build -o ./bin/ .

build-plugins:
    go build -o ./bin/plugins ./cmd/plugins

test-run:
    ./bin/fabric -path ./templates/ -plugins ./bin/plugins -document "test-document"

clean:
    rm -r ./bin/*

format:
    go mod tidy
    gofumpt -w .
    gci write --skip-generated -s standard -s default -s "prefix(github.com/blackstork-io/fabric)" .

format-extra: format
    gofumpt -w -extra .

lint: format
    go mod tidy
    gofumpt -w .
    gci write --skip-generated -s standard -s default -s "prefix(github.com/blackstork-io/fabric)" .
    golangci-lint run

test:
    go test -timeout 10s -race -v ./...