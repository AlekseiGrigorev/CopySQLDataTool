// Description: This package provides management features for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package app

import (
	"fmt"
	"os"
)

// FileProcessor represents a file processing utility that manages file operations.
// It contains a reference to an open file for writing or processing.
type FileProcessor struct {
	// Reference to an open file for writing or processing.
	File *os.File
}

// Write writes the given buffer of strings to the file, appending a newline
// after each statement. It does not use the data parameter.
// See: app.RowsProcessorInterface.Write
func (fp *FileProcessor) Write(buffer []string, data []any) error {
	for _, stmt := range buffer {
		_, err := fp.File.WriteString(stmt + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

// GetProcessedMsg returns a message indicating the number of rows processed
// to the file specified by the File field. This message is useful for logging
// and debugging purposes to confirm successful data processing.
// See: app.RowsProcessorInterface.GetProcessedMsg
func (fp *FileProcessor) GetProcessedMsg() string {
	return fmt.Sprint("Rows processed to file", fp.File.Name())
}
