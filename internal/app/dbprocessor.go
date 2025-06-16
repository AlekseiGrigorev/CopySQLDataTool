// Description: This package provides management features for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package app

import (
	"copysqldatatool/internal/appdb"
	"fmt"
	"strings"
)

// DbProcessor represents a database processor that manages database operations
// for a specific table using an AppDb connection.
type DbProcessor struct {
	// Database connection used for database operations.
	AppDb *appdb.AppDb
	// Name of the table to be used in database operations.
	TableName string
}

// Write executes a SQL statement with the given data on the database connection AppDb using the Exec method.
// The SQL statement is built by joining the strings in buffer with a space in between.
// The method returns an error if the execution of the SQL statement fails.
// See: app.RowsProcessorInterface.Write
func (db *DbProcessor) Write(buffer []string, data []any) error {
	_, err := db.AppDb.Exec(strings.Join(buffer, ""), data...)
	if err != nil {
		return fmt.Errorf("error writing to database: %w", err)
	}
	return nil
}

// GetProcessedMsg returns a message indicating the number of rows processed
// to the database table specified by the Table field. This message is useful
// for logging and debugging purposes to confirm successful data processing.
// See: app.RowsProcessorInterface.GetProcessedMsg
func (db *DbProcessor) GetProcessedMsg() string {
	return fmt.Sprint("Rows processed to table", db.TableName)
}
