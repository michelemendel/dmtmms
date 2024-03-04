# Member Management System


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
- https://hyperscript.org/docs/#install
  - Get the _hyperscript.min.js from https://github.com/bigskysoftware/_hyperscript/tree/master/dist
- sweetalert2 - 11.10.4
  - https://sweetalert2.github.io/
  - Download from
    - https://github.com/sweetalert2/sweetalert2/releases
      - sweetalert2.min.js
      - sweetalert2.min.css
  
Archived
- TW Elements (use Tailwind)
  - NOTE: I removed it, since it was more effective to use Tailwind directly.
  - https://tw-elements.com/
  - There were some issues: This didn't work (Uncaught SyntaxError: Unexpected token '<' (at tw-elements.umd.min.js:1:1))
  -	<script type="text/javascript" src="../node_modules/tw-elements/dist/js/tw-elements.umd.min.js"></script>
  - This worked, i.e. I had to get the file from the node_modules folder and put it in the public folder.
  - <script type="text/javascript" src="/public/tw-elements.umd.min.js"></script>			


### Production

Linux Debian 12

- Git
- Go (see go.mod for packages used)
- sqlite3 
- .bashrc
- dmtmms .env file

---

## Development environment

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

## Production environment

Templates and CSS are generated on the development machine and pushed to Github.

- Pull from Github
- $> make server
  - This will first build the application

- Use Github actions
  - Setup systemd - https://strapengine.com/auto-restart-mechanism-for-golang-program/
- Setup a cron job to backup the database

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

### Restore

- sqlite3 mydb.db ".restore '20240123T1658_mydb.db'"

### Some sqlite3 CLI commands

- .show
- .stats
- .tables
- .schema
- pragma busy_timeout = 5000
- explain query plan select ...
- .mode column , box, ...