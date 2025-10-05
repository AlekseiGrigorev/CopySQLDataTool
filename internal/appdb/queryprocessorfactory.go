// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

// QueryProcessorFactory is a factory for creating QueryProcessorInterface instances.
type QueryProcessorFactory struct{}

// CreateQueryProcessor creates a new instance of a query processor based on the specified query type.
// It initializes the query processor with the provided SQL query and sets the initial values using the
// given key-value map. The function returns an instance of QueryProcessorInterface that corresponds
// to the query type, supporting simple, limit-offset, and order-by-id query processing.
func (f *QueryProcessorFactory) CreateQueryProcessor(queryType string, query string, values map[string]any) QueryProcessorInterface {
	var p QueryProcessorInterface
	switch queryType {
	case QUERY_TYPE_LIMIT_OFFSET:
		p = &QueryProcessorLimitOffset{Query: query}
	case QUERY_TYPE_ORDERBYID:
		p = &QueryProcessorOrderByID{Query: query}
	case QUERY_TYPE_BETWEEN:
		p = &QueryProcessorBetween{Query: query}
	default:
		p = &QueryProcessorSimple{Query: query}
	}

	for key, value := range values {
		p.SetValue(key, value)
	}
	return p
}
