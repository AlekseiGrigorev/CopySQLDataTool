// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"fmt"
	"strings"
)

// QueryProcessorLimitOffset is a struct that implements the QueryProcessorInterface interface
// for processing SQL queries with LIMIT and OFFSET clauses for pagination.
type QueryProcessorLimitOffset struct {
	Query     string
	Limit     int64
	Offset    int64
	MaxOffset int64
}

// Return the type name for a simple query processor.
func (q *QueryProcessorLimitOffset) GetType() string {
	return QUERY_TYPE_LIMIT_OFFSET
}

// InitQuery resets the query processor to its initial state by setting the limit to 1000
// and the offset to 0.
func (q *QueryProcessorLimitOffset) InitQuery() QueryProcessorInterface {
	q.Limit = 1000
	q.Offset = 0
	q.MaxOffset = 0
	return q
}

// setValue sets the value of the specified key in the query processor.
func (q *QueryProcessorLimitOffset) SetValue(key string, value any) QueryProcessorInterface {
	switch strings.ToLower(key) {
	case "limit":
		q.Limit = value.(int64)
	case "offset":
		q.Offset = value.(int64)
	case "max_offset":
		q.MaxOffset = value.(int64)
	}
	return q
}

// ProcessQuery appends a LIMIT and OFFSET clause to the SQL query to enable pagination.
// It trims any trailing whitespace or semicolons from the original query before appending
// the LIMIT clause with the specified limit and offset values. The offset is incremented
// by the limit value after each call, facilitating paging through results.
// If MaxOffset is set, it ensures we don't exceed it.
func (q *QueryProcessorLimitOffset) ProcessQuery() string {
	trimmedQuery := strings.TrimRight(q.Query, " \t\n\r;")
	if q.MaxOffset > 0 && q.Offset > q.MaxOffset {
		query := trimmedQuery + " LIMIT 0 OFFSET 0;"
		return query
	}
	query := trimmedQuery + fmt.Sprintf(" LIMIT %d OFFSET %d;", q.Limit, q.Offset)
	q.Offset += q.Limit
	return query
}
