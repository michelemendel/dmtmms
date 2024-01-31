# DMT Member Management System


## Description

https://github.com/michelemendel/dmtmms

--- 

## Installation

### Development

Macos

- Git
- Go (see go.mod for packages used)
- sqlite3 
  - comes preinstalled with macOS
- Node
- Tailwind (CSS)
  - npm install -D tailwindcss
- HTMX - download htmx.min.js from https://htmx.org/
- templ 
  - https://github.com/a-h/templ
  - https://templ.guide/
  - go get github.com/a-h/templ
  - templ cli
    - go install github.com/a-h/templ/cmd/templ@latest

### Production

Linux Debian 12

- Git
- Go (see go.mod for packages used)
- sqlite3 
- .bashrc
- dmtmms .env file

---

## Start development environment

See also Makefile

- server (alternatives)
  - $> make server
  - $> make serverwatch
- template generation (alternatives)
  - $> make templ
  - $> make templwatch
- Tailwind (alternatives)
  - make tail
  - make tailwatch
- Start everything
  - $> make dev

Note: Tailwind is configured to look for HTML and JavaScript in: 
- public/ - html and js
- view/ - *_templ.go

See tailwind.config.js for more details

---

## Start production environment

Templates and CSS are generated on the development machine and pushed to Github.

- Pull from Github
- $> make server
  - This will first build the application

---

## CLI

This a command line interface, mainly used for database migrations.

- $> make cli

---

## Database: SQLITE3

This is a file and doesn't require a server

### Recommended

- pragma journal_mode = WAL
- pragma busy_timeout = 5000
- dates as ISO8601 strings, in Go it's RC3339
- STRICT mode

setting PRAGMA in Go, see https://gist.github.com/dgsb/6061941d2185f761848b143f080f4cd9

### Backup

- sqlite3 mydb.db ".backup '20240123T1658_mydb.db'"
  - This file will have to be moved to a safe place
- Alt. use litestream

### Some sqlite3 CLI commands

- .show
- .stats
- .tables
- .schema
- pragma busy_timeout = 5000
- explain query plan select ...
- .mode column , box, ...