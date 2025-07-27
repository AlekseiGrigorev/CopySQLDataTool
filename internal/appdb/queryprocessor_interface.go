// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

// QueryProcessorInterface is an interface for processing database queries.
type QueryProcessorInterface interface {
	// GetType returns the type name of the query processor.
	GetType() string

	// InitQuery resets the query processor to its initial state.
	InitQuery() QueryProcessorInterface

	// setValue sets the value of the specified key in the query processor.
	SetValue(key string, value any) QueryProcessorInterface

	// ProcessQuery processes the database query and returns the result as a string.
	ProcessQuery() string
}
