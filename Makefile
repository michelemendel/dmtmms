# --------------------------------------------------------------------------------
# Build

build_server:
	@go build -o bin/server ./cmd/api/...

set_cgo:
	go env -w CGO_ENABLED=1

# --------------------------------------------------------------------------------
# Run

templ:
	@echo "Generating templ files"
	@templ generate

templwatch:
	@echo "Generating and watching templ files"
	@templ generate --watch

server: build_server
	@echo "Starting web server"
	@./bin/server

serverwatch: build_server
	@echo "Starting web server with air to watch for changes"
	@air

tail:
	@echo "Generating CSS with Tailwind"
	@npx tailwindcss -i ./tailw_src/input.css -o ./public/output.css

# Also watches templ files
tailwatch:
	@echo "Generating and watching CSS with Tailwind"
	@npx tailwindcss -i ./tailw_src/input.css -o ./public/output.css --watch

# Starts and watches everything
watch_all:
	@${MAKE} -j3 templwatch serverwatch tailwatch

# --------------------------------------------------------------------------------
# CLI, mainly for database migrations

build_cli:
	@go build -o bin/cli ./cmd/cli/...

# --------------------------------------------------------------------------------
# Import data from a CSV file that was exported from the old member system
# This is a one-time job

build_import:
	@go build -o bin/import ./cmd/import/...

import: build_import
	@./bin/import