// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Constants for testing.
const (
	TEST_INSERT_INTO = "INSERT INTO "
	TEST_SELECT_FROM = "SELECT * FROM "
)

// prepareDb returns a DataReader instance configured to connect to a MySQL database.
// It initializes the DataReader with a default query that selects all columns from
// the test table and a limit of 1 row. The database connection is opened before
// returning the DataReader instance.
func prepareDr(t *testing.T) *DataReader {
	dr := DataReader{
		AppDb: &AppDb{Driver: "mysql", Dsn: TEST_DSN},
		Limit: 1,
		Query: TEST_SELECT_FROM_SQL,
	}
	err := dr.Open()
	if err != nil {
		t.Error(err)
	}
	return &dr
}

// prepareDbOrderById returns a DataReader instance configured to connect to a MySQL database.
// It initializes the DataReader with a default query that selects all columns from
// the test table where the id is greater than the InitialId field, ordered by id,
// and a limit of 1 row. The database connection is opened before
// returning the DataReader instance.
func prepareDrOrderById(t *testing.T) *DataReader {
	dr := DataReader{
		AppDb:         &AppDb{Driver: "mysql", Dsn: TEST_DSN},
		Limit:         1,
		Query:         TEST_SELECT_FROM + TEST_TBL_NAME + " WHERE id > {{id}} ORDER BY id LIMIT 1;",
		QueryType:     QUERY_TYPE_ORDERBYID,
		ExecutionTime: 3,
	}
	err := dr.Open()
	if err != nil {
		t.Error(err)
	}
	return &dr
}

// insertTestRows truncates the test table and inserts a specified number of rows into it.
// It executes SQL commands to insert rows with incremental integer values starting from 1 up to rowsCount.
// The function returns the modified DataReader instance. If any error occurs during insertion, it logs the error
// using the testing instance.
func insertTestRows(t *testing.T, dr *DataReader, rowsCount int) *DataReader {
	dr.AppDb.Exec(TEST_TRUNC_TBL_SQL)
	for i := 1; i <= rowsCount; i++ {
		_, err := dr.AppDb.Exec(TEST_INSERT_INTO + TEST_TBL_NAME + fmt.Sprintf(" VALUES (%d)", i))
		if err != nil {
			t.Error(err)
		}
	}
	return dr
}

// readTestRows reads rows from the DataReader and returns the number of rows read.
// It iterates over the result set, calling Next to advance the cursor and Scan to
// retrieve the row data. If an error occurs during iteration or scanning, the error
// is logged using the testing instance. The function returns the total count of rows
// successfully read and printed to the console.
func readTestRows(t *testing.T, dr *DataReader) int {
	counter := 0
	for {
		next, err := dr.Next()
		if err != nil {
			t.Error(err)
			break
		}
		if !next {
			break
		}
		counter++
		res, err := dr.Scan()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(res)
	}
	return counter
}

// TestDataReaderSimple tests the DataReader struct by verifying that it correctly reads all rows from a table.
// It initializes the DataReader with a default query that selects all columns from the test table, truncates the table,
// inserts 10 rows, and reads all rows. The test verifies that the correct number of rows is read by asserting that the
// counter equals 10.
func TestDataReaderSimple(t *testing.T) {
	dr := prepareDr(t)
	defer dr.Close()
	dr = insertTestRows(t, dr, 10)
	counter := readTestRows(t, dr)
	assert.Equal(t, counter, 10)
}

// TestDataReaderLimitOffset tests the DataReader struct by verifying that it correctly reads all rows from a table
// when configured with a limit and offset. It initializes the DataReader with a default query that selects all columns
// from the test table, truncates the table, inserts 10 rows, and reads all rows. The test verifies that the correct
// number of rows is read by asserting that the counter equals 10.
func TestDataReaderLimitOffset(t *testing.T) {
	dr := prepareDr(t)
	defer dr.Close()
	dr = insertTestRows(t, dr, 10)
	dr.QueryType = QUERY_TYPE_LIMIT_OFFSET
	dr.Limit = 4
	counter := readTestRows(t, dr)
	assert.Equal(t, counter, 10)
}

func TestDataReaderLimitOffsetInitialOffset(t *testing.T) {
	dr := prepareDr(t)
	defer dr.Close()
	dr = insertTestRows(t, dr, 10)
	dr.QueryType = QUERY_TYPE_LIMIT_OFFSET
	dr.Limit = 4
	dr.InitialOffset = 4
	counter := readTestRows(t, dr)
	assert.Equal(t, counter, 6)
}

func TestDataReaderLimitOffsetMaxOffset(t *testing.T) {
	dr := prepareDr(t)
	defer dr.Close()
	dr = insertTestRows(t, dr, 10)
	dr.QueryType = QUERY_TYPE_LIMIT_OFFSET
	dr.Limit = 4
	dr.MaxOffset = 4
	counter := readTestRows(t, dr)
	assert.Equal(t, counter, 8)
}

// TestDataReaderOrderById tests the DataReader struct by verifying that it correctly reads all rows from a table
// when configured with TYPE_ORDERBYID. It initializes the DataReader with a default query that selects all columns
// from the test table, truncates the table, inserts 10 rows, and reads all rows. The test verifies that the correct
// number of rows is read by asserting that the counter equals 10.
func TestDataReaderOrderById(t *testing.T) {
	dr := prepareDr(t)
	defer dr.Close()
	dr = insertTestRows(t, dr, 10)
	dr.QueryType = QUERY_TYPE_ORDERBYID
	dr.Query = "SELECT * FROM test_table WHERE id > {{id}} ORDER BY id LIMIT 4;"
	counter := readTestRows(t, dr)
	assert.Equal(t, counter, 10)
}

// TestDataReaderEmpty verifies that the DataReader returns 0 rows when the table is empty.
func TestDataReaderEmpty(t *testing.T) {
	dr := prepareDr(t)
	defer dr.Close()
	dr.AppDb.Exec(TEST_TRUNC_TBL_SQL)
	counter := readTestRows(t, dr)
	assert.Equal(t, counter, 0)
}

// TestDataReaderOne verifies that the DataReader correctly returns 1 row when the table contains a single entry.
// It initializes the database with one row and checks that the counter increments correctly to reflect the single row read.
func TestDataReaderOne(t *testing.T) {
	dr := prepareDr(t)
	defer dr.Close()
	dr.AppDb.Exec(TEST_TRUNC_TBL_SQL)
	dr.AppDb.Exec(TEST_INSERT_INTO + TEST_TBL_NAME + " VALUES (1)")
	counter := readTestRows(t, dr)
	assert.Equal(t, counter, 1)
}

// TestDataReaderMany verifies that the DataReader correctly returns all rows when the table contains multiple entries.
// It initializes the database with 3 rows and checks that the counter increments correctly to reflect all rows read.
func TestDataReaderMany(t *testing.T) {
	dr := prepareDr(t)
	defer dr.Close()
	dr.AppDb.Exec(TEST_TRUNC_TBL_SQL)
	dr.AppDb.Exec(TEST_INSERT_INTO_RAW_SQL)
	counter := readTestRows(t, dr)
	assert.Equal(t, counter, 3)
}

// TestDataReaderManyOrderById verifies that the DataReader correctly returns all rows in the correct order
// when the table contains multiple entries and the DataReader is configured with TYPE_ORDERBYID.
// It initializes the database with 3 rows and checks that the counter increments correctly to reflect all rows read.
func TestDataReaderManyOrderById(t *testing.T) {
	dr := prepareDrOrderById(t)
	defer dr.Close()
	dr.AppDb.Exec(TEST_TRUNC_TBL_SQL)
	dr.AppDb.Exec(TEST_INSERT_INTO_RAW_SQL)
	counter := readTestRows(t, dr)
	assert.Equal(t, counter, 3)
}

// TestDataReaderManyOrderByIdWrongType verifies that the DataReader correctly returns all rows in the correct order
// when the table contains multiple entries, the DataReader is configured with TYPE_ORDERBYID, and the query is not of type
// TYPE_ORDERBYID. It initializes the database with 3 rows and checks that the counter increments correctly to reflect
// all rows read.
func TestDataReaderManyOrderByIdWrongType(t *testing.T) {
	dr := prepareDrOrderById(t)
	defer dr.Close()
	dr.AppDb.Exec(TEST_TRUNC_TBL_SQL)
	dr.AppDb.Exec(TEST_INSERT_INTO_RAW_SQL)
	dr.Query = TEST_SELECT_FROM + TEST_TBL_NAME
	dr.ExecutionTime = 1
	counter := readTestRows(t, dr)
	assert.Equal(t, counter, 3)
}
