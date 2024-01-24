# --------------------------------------------------------------------------------
# Build

build_server:
	@go build -o bin/server ./cmd/api/...
	

# --------------------------------------------------------------------------------
# Run

# NOTE: Use 
# > make server
# or
# > air

# @APP_ENV=development ./bin/server 2>&1
server: build_server
	@./bin/server

templ_gen:
	@templ generate

templ_gen_watch:
	@templ generate --watch

tailw:
	@npx tailwindcss -i ./tailw_src/input.css -o ./public/output.css --watch

prod: templ_gen build_server
	@APP_ENV=production ./bin/server