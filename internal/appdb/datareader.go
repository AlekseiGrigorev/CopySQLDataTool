// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"database/sql"
	"fmt"
	"time"
)

// DataReader represents a database query reader with configurable parameters for executing and managing database queries.
// It supports features like query pagination, execution time limits, and dynamic query parameter management.
type DataReader struct {
	AppDb         *AppDb
	Query         string
	Args          []any
	QueryType     string
	Params        map[string]any
	ExecutionTime int64
	// Initial Id for query type "orderbyid"
	InitialId int64
	// Limit for query type "limitoffset"
	Limit int64
	// Initial Offset for query type "limitoffset"
	InitialOffset  int64
	queryProcessor QueryProcessorInterface
	columns        []string
	rows           *sql.Rows
	valuePtrs      []any
	values         []any
	startTime      time.Time
	prevQuery      string
	lastQuery      string
}

// Open opens the database connection for the underlying AppDb instance.
// It returns an error if the connection fails.
func (dataReader *DataReader) Open() error {
	return dataReader.AppDb.Open()
}

// Close closes the database connection for the underlying AppDb instance and resets all DataReader state fields.
// It is important to call Close after using a DataReader to ensure the database connection is released.
func (dataReader *DataReader) Close() {
	dataReader.closeRows()
	dataReader.startTime = time.Time{}
	dataReader.prevQuery = ""
	dataReader.lastQuery = ""
	dataReader.columns = nil
	dataReader.valuePtrs = nil
	dataReader.values = nil
	dataReader.AppDb.Close()
}

// closeRows closes the underlying sql.Rows if it is not nil. It is called by Close to ensure the
// underlying sql.Rows is released.
func (dataReader *DataReader) closeRows() {
	if dataReader.rows != nil {
		dataReader.rows.Close()
		dataReader.rows = nil
	}
}

// reopenAppDbByExecutionTime checks if the query execution time has exceeded the allowed
// ExecutionTime. If it has, it closes the underlying sql.Rows and the AppDb database
// connection, and then reopens the AppDb connection. It also resets the startTime to the
// current time. If the query execution time has not exceeded the allowed ExecutionTime,
// it does nothing.
func (dataReader *DataReader) reopenAppDbByExecutionTime() {
	if dataReader.ExecutionTime <= 0 {
		return
	}
	if dataReader.startTime.IsZero() {
		dataReader.startTime = time.Now()
		return
	}
	if time.Since(dataReader.startTime).Seconds() > float64(dataReader.ExecutionTime) {
		dataReader.closeRows()
		dataReader.AppDb.Close()
		dataReader.AppDb.Open()
		dataReader.startTime = time.Now()
		return
	}
}

// query executes the prepared SQL query and sets the rows and columns of the DataReader instance.
// It also handles connection reopening if the query execution time has exceeded the allowed
// ExecutionTime and handles the case when the AppDb connection is not open.
// It returns an error if the query execution fails.
func (dataReader *DataReader) query() error {
	if dataReader.queryProcessor == nil {
		dataReader.initQueryProcessor()
	}
	dataReader.queryProcessor.SetValue("id", dataReader.getLatId())
	query := dataReader.queryProcessor.ProcessQuery()
	dataReader.prevQuery = dataReader.lastQuery
	dataReader.lastQuery = query

	// Protection against error :Error 3024 (HY000): Query execution was interrupted, maximum statement execution time exceeded
	dataReader.reopenAppDbByExecutionTime()

	if !dataReader.AppDb.IsOpen() {
		err := dataReader.AppDb.Open()
		if err != nil {
			return err
		}
		defer dataReader.AppDb.Close()
	}

	dataReader.closeRows()

	rows, err := dataReader.AppDb.Query(query, dataReader.Args...)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	dataReader.columns = columns
	dataReader.rows = rows
	dataReader.valuePtrs = make([]any, len(columns))
	dataReader.values = make([]any, len(columns))
	for i := range dataReader.values {
		dataReader.valuePtrs[i] = &dataReader.values[i]
	}

	return nil
}

// initQueryProcessor initializes the query processor for the DataReader by creating a new
// instance using the QueryProcessorFactory. It sets the initial values for the query
// processor, including the initial ID, limit, and offset, and assigns the query processor
// to the DataReader's queryProcessor field.
func (dataReader *DataReader) initQueryProcessor() {
	queryProcessorFactory := QueryProcessorFactory{}
	values := make(map[string]any)
	values["id"] = dataReader.InitialId
	values["limit"] = dataReader.Limit
	values["offset"] = dataReader.InitialOffset
	dataReader.queryProcessor = queryProcessorFactory.CreateQueryProcessor(dataReader.QueryType, dataReader.Query, values)
}

// getLatId returns the last ID in the result set of the database query.
// If the query returned no rows, it returns the InitialId field.
func (dataReader *DataReader) getLatId() int64 {
	// If the values slice is empty or the first element is nil (no rows returned), return the initial ID
	if len(dataReader.values) == 0 || dataReader.values[0] == nil {
		return dataReader.InitialId
	}
	return dataReader.values[0].(int64)
}

// Columns returns a slice of strings containing the names of the columns
// in the result set of the database query.
func (dataReader *DataReader) Columns() []string {
	return dataReader.columns
}

// WrappedColumns returns a slice of strings containing the names of the columns
// in the result set of the database query, each wrapped in backticks.
// This is useful for formatting column names for SQL queries that require
// column names to be enclosed in backticks.
func (dataReader *DataReader) WrappedColumns() []string {
	cols := make([]string, len(dataReader.columns))
	for i, col := range dataReader.columns {
		cols[i] = fmt.Sprintf("`%s`", col)
	}
	return cols
}

// Next reads the next row from the database query. It returns true if there is a next row,
// false otherwise. If an error occurs while reading the row, it returns false and the error.
// If the query is exhausted, it closes the database reader and returns false and nil.
// It also handles the case where the query type has changed by re-executing the query
// and checking if there is a next row. It is idempotent and can be called multiple times.
func (dataReader *DataReader) Next() (bool, error) {
	if dataReader.rows == nil {
		err := dataReader.query()
		if err != nil {
			return false, err
		}
	}
	hasNext := dataReader.rows.Next()
	if !hasNext {
		err := dataReader.query()
		if err != nil {
			return false, err
		}
		// Protection against infinite loop if the query does not match the specified query type
		if dataReader.lastQuery != dataReader.prevQuery {
			hasNext = dataReader.rows.Next()
		}
	}
	if !hasNext {
		dataReader.Close()
		return false, nil
	}
	return true, nil
}

// Scan reads the next row from the database query and returns a slice of any values.
// It returns an error if the scan operation fails.
// It is idempotent and can be called multiple times.
func (dataReader *DataReader) Scan() ([]any, error) {
	if err := dataReader.rows.Scan(dataReader.valuePtrs...); err != nil {
		return nil, err
	}
	return dataReader.values, nil
}
