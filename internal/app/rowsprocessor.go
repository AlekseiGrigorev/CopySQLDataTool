// Description: This package provides management features for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package app

import (
	"copysqldatatool/internal/appbuffer"
	"copysqldatatool/internal/appdb"
	"copysqldatatool/internal/applog"
	"fmt"
)

// RowsProcessor manages the processing of database rows for data transfer or manipulation.
// It handles reading data, formatting, buffering, and writing rows with configurable processing.
type RowsProcessor struct {
	// Interface for processing rows.
	Processor RowsProcessorInterface
	// Data reader for retrieving rows from the source database.
	DataReader *appdb.DataReader
	// Log for recording processing details.
	Log *applog.AppLog
	// Dataset configuration for SQL insertion operations.
	Dataset Dataset
	// Buffer for storing formatted rows.
	buffer *appbuffer.AppBuffer
	// Data to be written to the processor.
	data []any
	// Formatter for formatting rows.
	formatter *appdb.Formatter
	// Columns to be used for formatting.
	columns []string
	// Count of rows processed in one insert command.
	count int
	// All processed rows counter.
	rowsCount int
}

// Process opens the data reader, reads rows, formats them according to the set InsertCommand and SqlStatement,
// and writes the formatted rows to the processor. It also handles closing the data reader and processing any remaining
// rows.
func (rp *RowsProcessor) Process() error {
	rp.reset()
	err := rp.DataReader.Open()
	if err != nil {
		return fmt.Errorf("error opening data reader: %w", err)
	}
	defer rp.DataReader.Close()

	for {
		next, err := rp.processRow()
		if err != nil {
			return err
		}
		if !next {
			break
		}
	}

	if rp.buffer.Len() > 0 {
		rp.buffer.AppendStr(";")
		if err := rp.Processor.Write(rp.buffer.GetBuffer(), rp.data); err != nil {
			return fmt.Errorf("error writing buffer to file: %w", err)
		}
		rp.buffer.Clear()
		rp.data = make([]any, 0)
	}
	if rp.Log != nil {
		rp.Log.Ok(rp.Processor.GetProcessedMsg(), ":", rp.rowsCount)
	}
	return nil
}

// reset resets the RowsProcessor to its initial state. It resets the count and rowsCount, clears the columns,
// resets the formatter, buffer, and data.
func (rp *RowsProcessor) reset() {
	rp.count = 0
	rp.rowsCount = 0
	rp.columns = make([]string, 0)
	rp.formatter = &appdb.Formatter{}
	rp.buffer = &appbuffer.AppBuffer{}
	rp.data = make([]any, 0)
}

// processRow reads the next row from the data reader, formats it according to the set SqlStatement,
// appends it to the buffer, and writes the buffer to the processor if the buffer is full.
// It also handles resetting the buffer and data if the buffer is full.
// Returns true if there is more data to be processed, false otherwise.
func (rp *RowsProcessor) processRow() (bool, error) {
	next, err := rp.DataReader.Next()
	if err != nil {
		return false, fmt.Errorf("error reading next row: %w", err)
	}

	if !next {
		return false, nil
	}

	if rp.rowsCount == 0 {
		rp.columns = rp.DataReader.WrappedColumns()
	}

	values, err := rp.DataReader.Scan()
	if err != nil {
		return false, fmt.Errorf("error scanning row: %w", err)
	}

	insertStatement := rp.formatter.GetInsertStatement(rp.Dataset.SqlStatementType, values)
	rp.appendRowToBuffer(insertStatement)
	if rp.Dataset.SqlStatementType == STATEMENT_TYPE_PREPARED {
		rp.data = append(rp.data, values...)
	}
	rp.count++
	rp.rowsCount++

	if rp.count == rp.Dataset.RowsPerCommand {
		rp.buffer.AppendStr(";")
		if err := rp.Processor.Write(rp.buffer.GetBuffer(), rp.data); err != nil {
			return false, fmt.Errorf("error writing buffer to file: %w", err)
		}
		rp.buffer.Clear()
		rp.data = make([]any, 0)
		rp.count = 0
		if rp.Log != nil {
			rp.Log.Info(rp.Processor.GetProcessedMsg(), "...:", rp.rowsCount)
		}
	}
	return true, nil
}

// appendRowToBuffer appends a row to the buffer in the correct format for the current SQL statement.
// If the buffer is empty, it adds the INSERT command and the first row in parentheses.
// If the buffer is not empty, it simply appends the next row in parentheses, separated by a comma.
func (rp *RowsProcessor) appendRowToBuffer(insertStatement string) {
	if rp.count == 0 {
		rp.buffer.AppendStr(rp.formatter.GetInsertCommand(rp.Dataset.InsertCommand, rp.Dataset.TableName, rp.columns))
		rp.buffer.AppendStr(fmt.Sprintf("(%s)", insertStatement))
	} else {
		rp.buffer.AppendStr(fmt.Sprintf(", (%s)", insertStatement))
	}
}
