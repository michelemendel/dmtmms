# DMT Member Management System

## Description


## Installation


## Usage


## SQLITE3

### Recommended

- pragma journal_mode = WAL
- pragma busy_timeout = 5000
- dates as ISO8601 strings, in Go it's RC3339

### Backup

- sqlite3 mydb.db ".backup '20240123T1658_mydb.db'"

### Look into

- setting PRAGMA in Go, see https://gist.github.com/dgsb/6061941d2185f761848b143f080f4cd9
- STRICT mode


create table if not exists ...
insert or ignore into...