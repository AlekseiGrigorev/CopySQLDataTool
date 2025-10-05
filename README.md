# CopySQLDataTool

## SQL data copy tool

Copy data from one database or instance to another. Works with MySQL, Clickhouse and another go database/sql interface databases.

* Flexible selective data copying based on queries - copy only what you need and save time and resources
* Copy data in parts - copy data in chunks to avoid overloading the database
* Multithreaded data copying from multiple tables - speed up the data copying process

## Repository

<https://github.com/AlekseiGrigorev/CopySQLDataTool>

## Build

Simple build
`go build -o copysqldatatool.exe main.go`
Build with static link
`go build -ldflags "-extldflags '-static'" -o copysqldatatool.exe main.go`

## Command line options

`-help` - show help (single option, without other options, show help and exit program)
`-version` - show version (single option, without other options, show version and exit program)
`-config <path to config file>` - path to config file (default: config.json)
`-log <path to log file>` - path to log file (default: no log file (on screen log))
`-go` - use goroutines (default: no (do not use goroutines))

For example: `./copysqldatatool.exe -config="config_local.json" -log="log.txt" -go`

## Config file

See file config.example.json

### Possible values

`$.config.source.driver, $.config.dest.driver` - DB driver name ("mysql", "clickhouse")

`$.config.default_dataset.copy_to, $.datasets.copy_to` - Copy data to ("file", "db" or "file,db")

`$.config.default_dataset.query_type, $.datasets.query_type` - Query type ("", "simple", "limitoffset", "orderbyid", "between")

For example:

Query type "simple" - `SELECT * FROM db.table`

Query type "limitoffset" - `SELECT * FROM db.table` - simple query that will be appended with the string `LIMIT %d OFFSET %d;` (where the LIMIT value is set, and OFFSET is calculated)

Query type "orderbyid" - `SELECT * FROM db.table WHERE id > {{id}} ORDER BY id LIMIT 10000;` - query with the placeholder `{{id}}` that will be replaced with the last id from the previous query

Query type "between" - `SELECT * FROM db.table WHERE field BETWEEN '{{start}}' AND '{{end}}' ...` - query with the placeholders `{{start}}` and `{{end}}` that will be replaced with the calculated values or dates from parameters

`$.config.default_dataset.sql_statement, $.datasets.sql_statement` - SQL statement ("prepared", "raw")

## Author

Aleksei Grigorev <https://www.aleksvgrig.com/>, <aleksvgrig@gmail.com>

## Copyright

Copyright (c) 2025 Aleksei Grigorev

## License

MIT License
