build_dir := "bin"
cli_binary_name := "shorty"
cli_main_package := "./cmd/shorty"
development_cli_config := "./config.yml"

# List available recipes
[group('help')]
default:
    @just --list

# Run all unit and integration tests
[group('tests')]
test:
    go test ./...

# Run all tests with Go's race detector
[group('tests')]
test-race:
    go test -race ./...

# Report statement coverage for all project packages
[group('tests')]
coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out

# Check that go.mod and go.sum are tidy without changing them
[group('quality')]
mod-tidy-check:
    go mod tidy -diff

# Run all repository quality checks
[group('quality')]
check: fmt-check mod-tidy-check vet test

# Format Go code
[group('quality')]
fmt:
    go fmt ./...

# Check that Go code is formatted
[group('quality')]
fmt-check:
    @files="$(gofmt -l $(rg --files -g '*.go'))"; if [ -n "$files" ]; then printf '%s\n' "$files"; exit 1; fi

# Run go vet
[group('quality')]
vet:
    go vet ./...

# Remove build artifacts
[group('artifacts')]
clean:
    rm -rf {{ build_dir }}

# Build the CLI and server binaries into ./bin
[group('artifacts')]
build: check
    mkdir -p {{ build_dir }}
    go build -o {{ build_dir }}/{{ cli_binary_name }} {{ cli_main_package }}

# Install the CLI and server binaries into GOPATH/bin or GOBIN
[group('artifacts')]
install: check
    go install {{ cli_main_package }}

# Run the CLI; pass arguments directly after `run`
[group('development')]
run *args:
    go run {{ cli_main_package }} --config {{ development_cli_config }} {{ quote(args) }}

# Run server subcommand; defaults to `serve`
[group('development')]
server:
    @just run serve
