// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"fmt"
	"strings"
)

// QueryProcessorOrderByID is a struct that implements the QueryProcessorInterface interface
// for processing SQL queries with the {{id}} placeholder replaced by the current value of the Id field.
type QueryProcessorOrderByID struct {
	Query string
	Id    int64
}

// Return the type name for a simple query processor.
func (q *QueryProcessorOrderByID) GetType() string {
	return QUERY_TYPE_ORDERBYID
}

// InitQuery resets the query processor to its initial state by setting the
// value of the Id field to 0.
func (q *QueryProcessorOrderByID) InitQuery() QueryProcessorInterface {
	q.Id = 0
	return q
}

// SetValue sets the value of the specified key in the query processor.
// The only key supported currently is "id", which is used to set the value of the
// {{id}} placeholder in the query string.
func (q *QueryProcessorOrderByID) SetValue(key string, value any) QueryProcessorInterface {
	switch strings.ToLower(key) {
	case "id":
		q.Id = value.(int64)
	}
	return q
}

// ProcessQuery implements the QueryProcessorInterface and returns the query string
// with the {{id}} placeholder replaced by the current value of the Id field.
func (q *QueryProcessorOrderByID) ProcessQuery() string {
	return strings.ReplaceAll(q.Query, "{{id}}", fmt.Sprintf("%d", q.Id))
}
