build_dir := "bin"
binary_name := "shorty"
main_package := "./cmd/shorty"
local_config := justfile_directory() + "/config.yml"

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

# Run all repository quality checks after regenerating database code
[group('quality')]
check: sqlc fmt-check mod-tidy-check vet test

# Format Go code
[group('quality')]
fmt:
    go fmt ./...

# Check that Go code is formatted
[group('quality')]
fmt-check:
    @files="$(gofmt -l $(rg --files -g '*.go'))"; if [ -n "$files" ]; then printf '%s\n' "$files"; exit 1; fi

# Run Go's static analysis
[group('quality')]
vet:
    go vet ./...

# Generate database code, run checks, and compile the single binary
[group('artifacts')]
build: check
    mkdir -p {{ build_dir }}
    go build -o {{ build_dir }}/{{ binary_name }} {{ main_package }}

# Generate database code, run checks, and install the single binary
[group('artifacts')]
install: check
    go install {{ main_package }}

# Generate database code and run Shorty using the local configuration
[group('development')]
run: sqlc
    SHORTY_CONFIG_FILE="{{ local_config }}" go run {{ main_package }}

# Generate database code and run Shorty with automatic rebuilds
[group('development')]
dev: sqlc
    SHORTY_CONFIG_FILE="{{ local_config }}" air -c .air.toml

# Generate the type-safe SQLite access layer from schema and queries
[group('database')]
sqlc:
    sqlc generate

# Create a new timestamped SQL migration: just migration add_something
[group('database')]
migration name:
    goose -dir internal/infra/sqlite/migrations create "{{ name }}" sql
