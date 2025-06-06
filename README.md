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

## Config file

See also config.example.json

```json
{
    "config": {
        "source": {
            "driver": "mysql",
            "dsn": "root:root@tcp(127.0.0.1:3306)/database"
        },
        "dest": {
            "driver": "clickhouse",
            "dsn": "clickhouse://default:default@127.0.0.1:9000/database2"
        },
        "default_dataset": {
            "insert_command": "INSERT INTO",
            "rows": 10000,
            "copy_to": "file,db",
            "query_type": "simple",
            "sql_statement": "prepared",
            "execution_time": 0
        }
    },
    "datasets": [
        {
            "query": "SELECT * FROM table",
            "table": "table",
            "enabled": true,
            "initial_id": 0
        },
        {
            "query": "SELECT * FROM table1",
            "table": "table1",
            "enabled": true,
            "initial_id": 0
        }
    ]
}
```

## Author

Aleksei Grigorev <https://www.aleksvgrig.com/>, <aleksvgrig@gmail.com>

## Copyright

Copyright (c) 2025 Aleksei Grigorev

## License

MIT License
