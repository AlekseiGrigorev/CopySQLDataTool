// Description: This package provides management features for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package app

// RowsProcessorInterface defines the contract for processing and writing rows of data
// with the ability to retrieve a processed message after completion.
// RowsProcessorInterface defines the contract for processing and writing rows of data
// with the ability to retrieve a processed message after completion.
type RowsProcessorInterface interface {
	// Write writes the provided buffer and data to the processor.
	// It returns an error if the write operation fails.
	Write(buffer []string, data []any) error

	// GetProcessedMsg retrieves the processed message after completion.
	// It returns the processed message as a string.
	GetProcessedMsg() string
}
