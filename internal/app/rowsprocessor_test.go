// Description: This package provides management features for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package app

import (
	"copysqldatatool/internal/appdb"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	FILE       = "test.txt"
	TBL_NAME_2 = "test_table2"
)

// prepareProcessor initializes a RowsProcessor with a DataReader that reads rows from a table in a MySQL database.
// It truncates the table, inserts 3 rows, and sets the DataReader to read the table with the given number of rows.
// The RowsProcessor is configured to write rows to the given processor as a single SQL statement.
// The function returns the prepared RowsProcessor or nil if there is an error.
func prepareProcessor(processor RowsProcessorInterface, rows int) *RowsProcessor {
	db := appdb.AppDb{
		Driver: "mysql",
		Dsn:    DSN,
	}
	err := db.Open()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	_, err = db.Exec(TRUNC_TBL_SQL)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	_, err = db.Exec(INSERT_3)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	p := RowsProcessor{
		Processor: processor,
		DataReader: &appdb.DataReader{
			AppDb: &db,
			Query: SELECT_TBL_SQL,
			Type:  appdb.TYPE_SIMPLE,
			Limit: rows,
		},
		Log: nil,
		Dataset: Dataset{
			RowsPerCommand:   rows,
			InsertCommand:    INSERT_INTO,
			TableName:        TBL_NAME_2,
			SqlStatementType: appdb.STATEMENT_RAW,
		},
	}

	return &p
}

// prepareDbRp returns a new AppDb instance with the database connection opened.
// It truncates the table given by TBL_NAME_2 after opening the database connection.
// It returns nil in case of any error during the process.
func prepareDbRp() *appdb.AppDb {
	db := appdb.AppDb{
		Driver: "mysql",
		Dsn:    DSN,
	}
	err := db.Open()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	_, err = db.Exec("TRUNCATE TABLE " + TBL_NAME_2)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &db
}

// cleanup removes the test file if it exists after a test is run.
func cleanup() {
	if _, err := os.Stat(FILE); err == nil {
		os.Remove(FILE)
	}
}

// TestWriteFileRp1 tests the Write method of the FileProcessor type.
// It creates a file, writes one row to it, and checks that the message returned
// by GetProcessedMsg contains the name of the file.
func TestWriteFileRp1(t *testing.T) {
	t.Cleanup(cleanup)
	file, err := os.Create(FILE)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, "error creating file")
	}
	defer file.Close()
	p := prepareProcessor(&FileProcessor{File: file}, 1)
	err = p.Process()
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, "error writing to file")
		return
	}
	assert.Empty(t, err)
}

// TestWriteFileRp2 tests the Write method of the FileProcessor type.
// It creates a file, writes two rows to it, and checks that the message returned
// by GetProcessedMsg contains the name of the file.
func TestWriteFileRp2(t *testing.T) {
	t.Cleanup(cleanup)
	file, err := os.Create(FILE)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, "error creating file")
	}
	defer file.Close()
	p := prepareProcessor(&FileProcessor{File: file}, 2)
	err = p.Process()
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, "error writing to file")
		return
	}
	assert.Empty(t, err)
}

// TestWriteDbRp tests the Write method of the DbProcessor type.
// It uses a DbProcessor that writes two rows to a MySQL database table
// using the STATEMENT_RAW SqlStatementType.
func TestWriteDbRp(t *testing.T) {
	p := prepareProcessor(&DbProcessor{AppDb: prepareDbRp(), TableName: TBL_NAME_2}, 2)
	p.Dataset.SqlStatementType = appdb.STATEMENT_RAW
	err := p.Process()
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, "error writing to db")
		return
	}
	assert.Empty(t, err)
}

// TestWriteDbRpPrepared tests the Write method of the DbProcessor type.
// It uses a DbProcessor that writes two rows to a MySQL database table
// using the STATEMENT_PREPARED SqlStatementType.
func TestWriteDbRpPrepared(t *testing.T) {
	p := prepareProcessor(&DbProcessor{AppDb: prepareDbRp(), TableName: TBL_NAME_2}, 2)
	p.Dataset.SqlStatementType = appdb.STATEMENT_PREPARED
	err := p.Process()
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, "error writing to db")
		return
	}
	assert.Empty(t, err)
}
