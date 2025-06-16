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

// AppDb represents a database connection configuration and handle.
// It contains the database driver, data source name (DSN), and an underlying SQL database connection.
type AppDb struct {
	// Driver is the database driver name.
	// It should be one of the supported database drivers.
	// For example, "mysql", "postgres", "sqlite3", etc.
	Driver string
	// Dsn is the data source name for the database connection.
	Dsn string
	// db is the underlying SQL database connection.
	db *sql.DB
}

// Open opens a database connection with the database driver and data source name (DSN)
// specified in the AppDb instance. If the connection is already open, it does nothing.
// It returns an error if the connection fails.
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

// Close closes the database connection if it is open.
// It sets the underlying SQL database connection to nil after closing.
// If the connection is already closed, it does nothing.
// Returns an error if the operation fails.
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

// IsOpen checks if the database connection is currently open.
// It returns true if the connection is open, otherwise it returns false.
func (appdb *AppDb) IsOpen() bool {
	return appdb.db != nil
}

// Exec executes a SQL statement with the given data on the database connection.
// The method returns a Result instance if the execution is successful, otherwise it returns an error.
// The Result instance provides information about the number of affected rows and the last inserted ID.
func (appdb *AppDb) Exec(sqlCommand string, args ...any) (result sql.Result, err error) {
	result, err = appdb.db.Exec(sqlCommand, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ExecMultiple executes multiple SQL statements on the database connection.
// The method takes a single string argument containing multiple SQL statements
// separated by semicolons. It splits the string into individual statements,
// trims any whitespace around each statement, and executes each statement
// using the Exec method. If any statement results in an error, it returns
// the error. Otherwise it returns nil.
func (appdb *AppDb) ExecMultiple(sqlCommands string) error {
	commands := strings.SplitSeq(sqlCommands, ";")
	for command := range commands {
		trimmedCommand := strings.TrimSpace(command)
		if trimmedCommand == "" {
			continue
		}
		_, err := appdb.db.Exec(command)
		if err != nil {
			return err
		}
	}
	return nil
}

// PrepareExec prepares a SQL statement and executes it with the given data on the database connection.
// The method takes a SQL statement string and a variable number of arguments.
// It prepares the statement using the database connection's Prepare method,
// executes it with the given arguments using the statement's Exec method,
// and returns the result if the execution is successful, otherwise it returns an error.
// The result is a sql.Result instance that provides information about the number of affected rows and the last inserted ID.
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

// QueryRow executes a SQL query with the given data on the database connection
// and returns the first row of the result set. If the query returns no rows, it
// returns a nil *sql.Row and a nil error. Otherwise it returns a *sql.Row
// instance that can be used to retrieve the columns of the row, and a nil error.
func (appdb *AppDb) QueryRow(sqlCommand string, args ...any) (result *sql.Row, err error) {
	return appdb.db.QueryRow(sqlCommand, args...), nil
}

// Query executes a SQL query with the given data on the database connection
// and returns the result set. If the query returns no rows, it returns a nil
// *sql.Rows and a nil error. Otherwise it returns a *sql.Rows instance that
// can be used to retrieve the columns and rows of the result set, and a nil
// error.
func (appdb *AppDb) Query(sqlCommand string, args ...any) (result *sql.Rows, err error) {
	return appdb.db.Query(sqlCommand, args...)
}

// GetScalar executes a SQL query with the given data on the database connection
// and returns the single value (scalar) of the first column of the first row of the result set.
// If the query returns no rows, it returns a nil value and a nil error.
// If the query returns multiple rows or columns, it only returns the first column of the first row.
// If the query returns an error, it returns a nil value and the error.
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
