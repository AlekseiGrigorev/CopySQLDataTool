// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"fmt"
	"strconv"
	"strings"
)

// Formatter provides methods for formatting database-related operations like insert statements and value formatting.
type Formatter struct {
}

// AppendInitialInsert appends an initial SQL INSERT command to the buffer.
// It constructs the command using the provided SQL command, table name, and columns,
// followed by the initial insert statement in parentheses. The formatted command is
// appended to the buffer, which is then returned.
func (f *Formatter) AppendInitialInsert(buffer []string, command string, table string, columns []string, insertStatement string) []string {
	columnsStr := strings.Join(columns, ", ")
	insertCommand := fmt.Sprintf("%s %s (%s) VALUES", command, table, columnsStr)
	buffer = append(buffer, insertCommand)
	buffer = append(buffer, fmt.Sprintf("(%s)", insertStatement))
	return buffer
}

// GetInsertCommand constructs and returns an SQL INSERT command string.
// The command includes the specified SQL command, table name, and columns,
// formatted as "COMMAND TABLE (columns) VALUES".
func (f *Formatter) GetInsertCommand(command string, table string, columns []string) string {
	columnsStr := strings.Join(columns, ", ")
	return fmt.Sprintf("%s %s (%s) VALUES", command, table, columnsStr)
}

// GetInsertStatement constructs and returns an SQL INSERT statement string.
// It takes the statement type and values as parameters.
// If the statement type is STATEMENT_PREPARED, it builds placeholders for the number of values provided.
// Otherwise, it formats the values according to the database type and returns the formatted string.
func (f *Formatter) GetInsertStatement(statement string, values []any) string {
	if statement == STATEMENT_TYPE_PREPARED {
		return f.BuildInsertPlaceholders(len(values))
	}
	return f.FormatRowValues(values)
}

// FormatValue formats a given value for SQL insertion based on its type.
// - For byte slices and strings, it escapes single quotes and backslashes and wraps the value in single quotes.
// - For integer types, it converts the value to a string representation of the number.
// - For float types, it converts the value to a string representation with no unnecessary precision.
// - For nil, it returns the SQL NULL keyword.
// - For all other types, it uses the default string representation wrapped in single quotes.
func (f *Formatter) FormatValue(val any) string {
	switch v := val.(type) {
	case []byte:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(strings.ReplaceAll(string(v), "'", "''"), "\\", "\\\\"))
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(strings.ReplaceAll(v, "'", "''"), "\\", "\\\\"))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// FormatRowValues formats a slice of values for SQL insertion.
// It iterates over the values, formats each one according to its type using FormatValue,
// and joins the formatted values with a comma separator.
func (f *Formatter) FormatRowValues(values []any) string {
	var formattedValues []string

	for _, val := range values {
		formattedValues = append(formattedValues, f.FormatValue(val))
	}

	return strings.Join(formattedValues, ", ")
}

// BuildInsertPlaceholders builds and returns a string of placeholders for a SQL INSERT statement.
// It takes the number of columns as a parameter and returns a string of the form "?, ?, ..., ?".
func (f *Formatter) BuildInsertPlaceholders(columnCount int) string {
	return strings.Repeat("?, ", columnCount-1) + "?"
}
