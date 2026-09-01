build_dir := "bin"
api_binary_name := "shorty-api"
api_main_package := "./cmd/api"
web_binary_name := "shorty-web"
web_main_package := "./cmd/web"
local_api_config := justfile_directory() + "/config-api.yml"
local_web_config := justfile_directory() + "/config-web.yml"
tailwind_input := "internal/web/styles/global.css"
tailwind_output_dir := "internal/web/static/css"
tailwind_output := "internal/web/static/css/style.css"

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

# Compile Tailwind and the project binaries into ./bin
[group('artifacts')]
build: css check
    mkdir -p {{ build_dir }}
    go build -o {{ build_dir }}/{{ api_binary_name }} {{ api_main_package }}
    go build -o {{ build_dir }}/{{ web_binary_name }} {{ web_main_package }}

# Install the API and Web binaries into GOPATH/bin or GOBIN
[group('artifacts')]
install: css check
    go install {{ api_main_package }}
    go install {{ web_main_package }}

# Run the JSON API using its configuration file
[group('development')]
api:
    SHORTY_API_CONFIG_FILE="{{ local_api_config }}" go run {{ api_main_package }}

# Run the server-rendered Web application using its configuration file
[group('development')]
web: css
    SHORTY_WEB_CONFIG_FILE="{{ local_web_config }}" go run {{ web_main_package }}

# Run API and Web together with automatic rebuilds
[group('development')]
[parallel]
dev: api-watch web-watch

# Rebuild and restart the JSON API when its source or configuration changes
[group('development')]
api-watch:
    SHORTY_API_CONFIG_FILE="{{ local_api_config }}" air -c .air.api.toml

# Rebuild Web and Tailwind, then reload the browser through localhost:3001
[group('development')]
web-watch:
    SHORTY_WEB_CONFIG_FILE="{{ local_web_config }}" air -c .air.web.toml

# Compile the minified Tailwind stylesheet used by the embedded Web UI
[group('styles')]
css:
    mkdir -p {{ tailwind_output_dir }}
    tailwindcss -i {{ tailwind_input }} -o {{ tailwind_output }} --minify

# Recompile Tailwind whenever templates or source styles change
[group('styles')]
css-watch:
    tailwindcss -i {{ tailwind_input }} -o {{ tailwind_output }} --watch
