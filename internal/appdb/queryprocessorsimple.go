// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

// QueryProcessorSimple is a simple implementation of the QueryProcessorInterface
// that returns the query string as is without any changes.
type QueryProcessorSimple struct {
	Query string
}

// Return the type name for a simple query processor.
func (q *QueryProcessorSimple) GetType() string {
	return QUERY_TYPE_SIMPLE
}

// InitQuery implements the QueryProcessorInterface and does nothing, as a simple query
// processor does not have any state that needs to be initialized.
// No need to init anything for a simple query processor.
func (q *QueryProcessorSimple) InitQuery() QueryProcessorInterface {
	return q
}

// setValue implements the QueryProcessorInterface and does nothing, as a simple query
// processor does not have any state that needs to be set.
func (q *QueryProcessorSimple) SetValue(key string, value any) QueryProcessorInterface {
	// No need to set any values for a simple query processor.
	return q
}

// ProcessQuery implements the QueryProcessorInterface and returns the query string
// as is without any changes.
func (q *QueryProcessorSimple) ProcessQuery() string {
	return q.Query
}
