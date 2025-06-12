// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type DataReader struct {
	AppDb         *AppDb
	Limit         int
	Query         string
	Args          []any
	Type          string
	Params        map[string]any
	ExecutionTime int
	InitialId     int
	nextOffset    int
	columns       []string
	rows          *sql.Rows
	valuePtrs     []any
	values        []any
	startTime     time.Time
	prevQuery     string
	lastQuery     string
}

func (dataReader *DataReader) Open() error {
	return dataReader.AppDb.Open()
}

func (dataReader *DataReader) Close() {
	dataReader.closeRows()
	dataReader.nextOffset = 0
	dataReader.startTime = time.Time{}
	dataReader.prevQuery = ""
	dataReader.lastQuery = ""
	dataReader.columns = make([]string, 0)
	dataReader.valuePtrs = make([]any, 0)
	dataReader.values = make([]any, 0)
	dataReader.AppDb.Close()
}

func (dataReader *DataReader) closeRows() {
	if dataReader.rows != nil {
		dataReader.rows.Close()
		dataReader.rows = nil
	}
}

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

func (dataReader *DataReader) prepareQueryLimitOffset() string {
	if dataReader.rows == nil {
		dataReader.nextOffset = 0
	}
	trimmedQuery := strings.TrimRight(dataReader.Query, " \t\n\r;")
	query := trimmedQuery + fmt.Sprintf(" LIMIT %d OFFSET %d;", dataReader.Limit, dataReader.nextOffset)
	dataReader.nextOffset += dataReader.Limit
	return query
}

func (dataReader *DataReader) prepareQueryOrderById() string {
	query := ""
	if len(dataReader.values) == 0 {
		query = strings.ReplaceAll(dataReader.Query, "{{id}}", fmt.Sprintf("%d", dataReader.InitialId))
	} else {
		query = strings.ReplaceAll(dataReader.Query, "{{id}}", fmt.Sprintf("%d", dataReader.values[0]))
	}
	return query
}

func (dataReader *DataReader) prepareQuery() string {
	switch dataReader.Type {
	case "limitoffset":
		return dataReader.prepareQueryLimitOffset()
	case "orderbyid":
		return dataReader.prepareQueryOrderById()
	}
	return dataReader.Query
}

func (dataReader *DataReader) query() error {
	query := dataReader.prepareQuery()
	dataReader.prevQuery = dataReader.lastQuery
	dataReader.lastQuery = query

	// Защита от ошибки :Error 3024 (HY000): Query execution was interrupted, maximum statement execution time exceeded
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

func (dataReader *DataReader) Columns() []string {
	return dataReader.columns
}

func (dataReader *DataReader) WrappedColumns() []string {
	cols := make([]string, len(dataReader.columns))
	for i, col := range dataReader.columns {
		cols[i] = fmt.Sprintf("`%s`", col)
	}
	return cols
}

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
		// Защита от бесконечного цикла если запрос не соответствует указанному типу запроса
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

func (dataReader *DataReader) Scan() ([]any, error) {
	if err := dataReader.rows.Scan(dataReader.valuePtrs...); err != nil {
		return nil, err
	}
	return dataReader.values, nil
}
