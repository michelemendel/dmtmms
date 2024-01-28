# --------------------------------------------------------------------------------
# Build

build_server:
	@go build -o bin/server ./cmd/api/...
	

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
# Database

build_migration:
	@go build -o bin/migration ./cmd/migration/...

migrate: build_migration
	@echo "Migrating database"
	@./bin/migration
