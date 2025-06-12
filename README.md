# CopySQLDataTool

## SQL data copy tool

Copy data from one database or instance to another. Works with MySQL, Clickhouse and another go database/sql interface databases.

## Repository

<https://github.com/AlekseiGrigorev/CopySQLDataTool>

## Build

Simple build
`go build -o copysqldatatool.exe main.go`
Build with static link
`go build -ldflags "-extldflags '-static'" -o copysqldatatool.exe main.go`

## Command line options

`-config <path to config file>` - path to config file (default: config.json)
`-log <path to log file>` - path to log file (default: no log file (on screen log))

## Config file

See file config.example.json

## Author

Aleksei Grigorev <https://www.aleksvgrig.com/>, <aleksvgrig@gmail.com>

## Copyright

Copyright (c) 2025 Aleksei Grigorev

## License

MIT License
