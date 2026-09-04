build_dir := "bin"
api_binary_name := "shorty-api"
api_main_package := "./cmd/api"
cli_binary_name := "shorty"
cli_main_package := "./cmd/cli"
api_config := justfile_directory() + "/config-api.yml"
cli_config := justfile_directory() + "/config-cli.yml"
tailwind_input := "internal/web/styles/global.css"
tailwind_output := "internal/web/static/css/style.css"
htmx_input := "node_modules/htmx.org/dist/htmx.min.js"
htmx_output_dir := "internal/web/static/js"
htmx_output := htmx_output_dir + "/htmx.min.js"
data_dir := justfile_directory() + "/data"

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

# Run all repository quality checks after regenerating code and Web assets
[group('quality')]
check: sqlc templ css htmx fmt-check mod-tidy-check vet test

# Format Go code
[group('quality')]
fmt:
    go fmt ./...
    go tool templ fmt .

# Check that Go code is formatted
[group('quality')]
fmt-check:
    @files="$(gofmt -l $(rg --files -g '*.go'))"; if [ -n "$files" ]; then printf '%s\n' "$files"; exit 1; fi
    go tool templ fmt -fail .

# Run Go's static analysis
[group('quality')]
vet:
    go vet ./...

# Remove build artifacts and reset the local SQLite database
[group('artifacts')]
clean:
    rm -rf "{{ build_dir }}"
    rm -rf "{{ data_dir }}"

# Generate code, run checks, and compile the API and CLI binaries
[group('artifacts')]
build: check
    mkdir -p {{ build_dir }}
    go build -o {{ build_dir }}/{{ api_binary_name }} {{ api_main_package }}
    go build -o {{ build_dir }}/{{ cli_binary_name }} {{ cli_main_package }}

# Run the Shorty CLI and forward its arguments
[group('development')]
run *args:
    SHORTY_CONFIG_FILE="{{ cli_config }}" go run {{ cli_main_package }} {{ args }}

# Generate code and assets, then serve the API with automatic rebuilds
[group('development')]
serve: sqlc templ css htmx
    SHORTY_CONFIG_FILE="{{ api_config }}" air -c .air.toml

# Compile the Web stylesheet with Tailwind CSS
[group('assets')]
css:
    pnpm exec tailwindcss -i {{ tailwind_input }} -o {{ tailwind_output }} --minify

# Copy the pinned HTMX browser bundle into the embedded static files
[group('assets')]
htmx:
    mkdir -p {{ htmx_output_dir }}
    cp {{ htmx_input }} {{ htmx_output }}

# Generate Go components from templ sources
[group('generation')]
templ:
    go tool templ generate

# Generate the type-safe SQLite access layer from schema and queries
[group('generation')]
sqlc:
    sqlc generate

# Create a new timestamped SQL migration: just migration add_something
[group('database')]
migration name:
    goose -dir internal/infra/sqlite/migrations create "{{ name }}" sql
