package app

import (
	"copysqldatatool/internal/appbuffer"
	"copysqldatatool/internal/appdb"
	"copysqldatatool/internal/applog"
	"fmt"
)

type RowsProcessor struct {
	Processor     RowsProcessorInterface
	DataReader    *appdb.DataReader
	Log           *applog.AppLog
	InsertCommand string
	Table         string
	Rows          int
	SqlStatement  string
	buffer        *appbuffer.AppBuffer
	data          []any
	formatter     *appdb.Formatter
	columns       []string
	count         int
	rowsCount     int
}

func (rp *RowsProcessor) init() {
	rp.count = 0
	rp.rowsCount = 0
	rp.columns = make([]string, 0)
	rp.formatter = &appdb.Formatter{}
	rp.buffer = &appbuffer.AppBuffer{}
	rp.data = make([]any, 0)
}

func (rp *RowsProcessor) Process() error {
	rp.init()
	err := rp.DataReader.Open()
	if err != nil {
		return fmt.Errorf("error opening data reader: %w", err)
	}
	defer rp.DataReader.Close()

	for {
		next, err := rp.processRow()
		if err != nil {
			return (err)
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

	insertStatement := rp.formatter.GetInsertStatement(rp.SqlStatement, values)
	rp.appendRowToBuffer(insertStatement)
	if rp.SqlStatement == appdb.STATEMENT_PREPARED {
		rp.data = append(rp.data, values...)
	}
	rp.count++
	rp.rowsCount++

	if rp.count == rp.Rows {
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

func (rp *RowsProcessor) appendRowToBuffer(insertStatement string) {
	if rp.count == 0 {
		rp.buffer.AppendStr(rp.formatter.GetInsertCommand(rp.InsertCommand, rp.Table, rp.columns))
		rp.buffer.AppendStr(fmt.Sprintf("(%s)", insertStatement))
	} else {
		rp.buffer.AppendStr(fmt.Sprintf(", (%s)", insertStatement))
	}
}
