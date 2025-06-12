// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"database/sql"
	"strings"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/go-sql-driver/mysql"
)

const (
	TYPE_UNDEFINED     = ""
	TYPE_SIMPLE        = "simple"
	TYPE_LIMIT_OFFSET  = "limitoffset"
	TYPE_ORDERBYID     = "orderbyid"
	STATEMENT_PREPARED = "prepared"
	STATEMENT_RAW      = "raw"
)

type AppDb struct {
	Driver string
	Dsn    string
	db     *sql.DB
}

func (appdb *AppDb) Open() error {
	if appdb.db != nil {
		return nil
	}
	db, err := sql.Open(appdb.Driver, appdb.Dsn)
	if err != nil {
		return err
	}
	appdb.db = db
	return nil
}

func (appdb *AppDb) Close() error {
	if appdb.db == nil {
		return nil
	}
	err := appdb.db.Close()
	if err != nil {
		return err
	}
	appdb.db = nil
	return nil
}

func (appdb *AppDb) IsOpen() bool {
	return appdb.db != nil
}

func (appdb *AppDb) Exec(sqlCommand string, args ...any) (result sql.Result, err error) {
	result, err = appdb.db.Exec(sqlCommand, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (appdb *AppDb) ExecMultiple(sqlCommands string) (err error) {
	commands := strings.SplitSeq(sqlCommands, ";")
	for command := range commands {
		trimmedCommand := strings.TrimSpace(command)
		if trimmedCommand == "" {
			continue
		}
		_, err = appdb.db.Exec(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func (appdb *AppDb) PrepareExec(sqlCommand string, args ...any) (result sql.Result, err error) {
	stmt, err := appdb.db.Prepare(sqlCommand)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err = stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (appdb *AppDb) QueryRow(sqlCommand string, args ...any) (result *sql.Row, err error) {
	return appdb.db.QueryRow(sqlCommand, args...), nil
}

func (appdb *AppDb) Query(sqlCommand string, args ...any) (result *sql.Rows, err error) {
	return appdb.db.Query(sqlCommand, args...)
}

func (appdb *AppDb) GetScalar(sqlCommand string, args ...any) (result any, err error) {
	var value any
	res, err := appdb.QueryRow(sqlCommand, args...)
	if err != nil {
		return nil, err
	}
	err = res.Scan(&value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
