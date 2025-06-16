// Description: This package provides management features for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package app

const (
	// Statement type for prepared statements (INSERT INTO ... VALUES (?, ?, ...))
	STATEMENT_TYPE_PREPARED = "prepared"
	// Statement type for raw SQL statements (INSERT INTO ... VALUES ('value1', 'value2', ...))
	STATEMENT_TYPE_RAW = "raw"
)

// Dataset represents a database dataset configuration with details for SQL insertion operations.
// It contains parameters for constructing and executing SQL insert statements.
type Dataset struct {
	// InsertCommand is the part of SQL command used for inserting data into the database.
	// INSERT INTO or INSERT IGNORE INTO etc.
	InsertCommand string
	// Table name
	TableName string
	// Rows per command to insert multiple rows at once
	RowsPerCommand int
	// Type of SQL statement to be used for the insert operation.
	// prepared, simple, custom etc.
	// See: STATEMENT_TYPE_PREPARED, STATEMENT_TYPE_RAW
	SqlStatementType string
}
