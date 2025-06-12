// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"fmt"
	"strconv"
	"strings"
)

type Formatter struct {
}

func (f *Formatter) AppendInitialInsert(buffer []string, command string, table string, columns []string, insertStatement string) []string {
	columnsStr := strings.Join(columns, ", ")
	insertCommand := fmt.Sprintf("%s %s (%s) VALUES", command, table, columnsStr)
	buffer = append(buffer, insertCommand)
	buffer = append(buffer, fmt.Sprintf("(%s)", insertStatement))
	return buffer
}

func (f *Formatter) GetInsertCommand(command string, table string, columns []string) string {
	columnsStr := strings.Join(columns, ", ")
	return fmt.Sprintf("%s %s (%s) VALUES", command, table, columnsStr)
}

func (f *Formatter) GetInsertStatement(statement string, values []any) string {
	if statement == STATEMENT_PREPARED {
		return f.BuildInsertPlaceholders(len(values))
	}
	return f.FormatRowValues(values)
}

func (f *Formatter) FormatValue(val any) string {
	switch v := val.(type) {
	case []byte:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(strings.ReplaceAll(string(v), "'", "''"), "\\", "\\\\"))
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(strings.ReplaceAll(string(v), "'", "''"), "\\", "\\\\"))
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

func (f *Formatter) FormatRowValues(values []any) string {
	var formattedValues []string

	for _, val := range values {
		formattedValues = append(formattedValues, f.FormatValue(val))
	}

	return strings.Join(formattedValues, ", ")
}

func (f *Formatter) BuildInsertPlaceholders(columnCount int) string {
	return strings.Repeat("?, ", columnCount-1) + "?"
}
