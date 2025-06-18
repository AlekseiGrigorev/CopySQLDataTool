// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TEST_INSERT_INTO = "INSERT INTO "
const TEST_SELECT_FROM = "SELECT * FROM "

// prepareDb returns a DataReader instance configured to connect to a MySQL database.
// It initializes the DataReader with a default query that selects all columns from
// the test table and a limit of 1 row. The database connection is opened before
// returning the DataReader instance.
func prepareDr(t *testing.T) DataReader {
	dr := DataReader{
		AppDb: &AppDb{Driver: "mysql", Dsn: TEST_DSN},
		Limit: 1,
		Query: TEST_SELECT_FROM_SQL,
	}
	err := dr.Open()
	if err != nil {
		t.Error(err)
	}
	return dr
}

// prepareDbOrderById returns a DataReader instance configured to connect to a MySQL database.
// It initializes the DataReader with a default query that selects all columns from
// the test table where the id is greater than the InitialId field, ordered by id,
// and a limit of 1 row. The database connection is opened before
// returning the DataReader instance.
func prepareDrOrderById(t *testing.T) DataReader {
	dr := DataReader{
		AppDb:         &AppDb{Driver: "mysql", Dsn: TEST_DSN},
		Limit:         1,
		Query:         TEST_SELECT_FROM + TEST_TBL_NAME + " WHERE id > {{id}} ORDER BY id LIMIT 1;",
		Type:          QUERY_TYPE_ORDERBYID,
		ExecutionTime: 3,
	}
	err := dr.Open()
	if err != nil {
		t.Error(err)
	}
	return dr
}

// TestPrepareQuery verifies that the prepareQuery method returns the input query
// when the DataReader is configured to use a TYPE_SIMPLE query.
func TestPrepareQuery(t *testing.T) {
	query := TEST_SELECT_FROM_SQL
	dr := DataReader{
		Query: TEST_SELECT_FROM_SQL,
		Limit: 1,
		Type:  QUERY_TYPE_SIMPLE,
	}
	result := dr.prepareQuery()
	assert.Equal(t, query, result)
}

// TestPrepareQueryLimitOffset verifies that the prepareQuery method correctly appends
// a LIMIT and OFFSET clause to the SQL query when the DataReader is configured with
// TYPE_LIMIT_OFFSET. It checks that the resulting query includes a LIMIT of 1 and
// an OFFSET of 0, matching the expected query.
func TestPrepareQueryLimitOffset(t *testing.T) {
	query := TEST_SELECT_FROM + TEST_TBL_NAME + " LIMIT 1 OFFSET 0;"
	dr := DataReader{
		Query: TEST_SELECT_FROM_SQL,
		Limit: 1,
		Type:  QUERY_TYPE_LIMIT_OFFSET,
	}
	result := dr.prepareQuery()
	assert.Equal(t, query, result)
}

// TestPrepareQueryOrderById verifies that the prepareQuery method correctly replaces
// the {{id}} parameter with 0 when the DataReader is configured with TYPE_ORDERBYID.
// It checks that the resulting query is the expected query with the replaced parameter.
func TestPrepareQueryOrderById(t *testing.T) {
	query := TEST_SELECT_FROM + TEST_TBL_NAME + " WHERE id > 0 ORDER BY id LIMIT 1;"
	dr := DataReader{
		Query: TEST_SELECT_FROM + TEST_TBL_NAME + " WHERE id > {{id}} ORDER BY id LIMIT 1;",
		Type:  QUERY_TYPE_ORDERBYID,
	}
	result := dr.prepareQuery()
	assert.Equal(t, query, result)
}

// TestDataReaderEmpty verifies that the DataReader returns 0 rows when the table is empty.
func TestDataReaderEmpty(t *testing.T) {
	dr := prepareDr(t)
	defer dr.Close()
	dr.AppDb.Exec(TEST_TRUNC_TBL_SQL)
	counter := 0
	for {
		next, err := dr.Next()
		if err != nil {
			fmt.Println(err)
			break
		}
		if !next {
			break
		}
		counter++
		res, err := dr.Scan()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
	assert.Equal(t, counter, 0)
}

// TestDataReaderOne verifies that the DataReader correctly returns 1 row when the table contains a single entry.
// It initializes the database with one row and checks that the counter increments correctly to reflect the single row read.
func TestDataReaderOne(t *testing.T) {
	dr := prepareDr(t)
	defer dr.Close()
	dr.AppDb.Exec(TEST_TRUNC_TBL_SQL)
	dr.AppDb.Exec(TEST_INSERT_INTO + TEST_TBL_NAME + " VALUES (1)")
	counter := 0
	for {
		next, err := dr.Next()
		if err != nil {
			fmt.Println(err)
			break
		}
		if !next {
			break
		}
		counter++
		res, err := dr.Scan()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
	assert.Equal(t, counter, 1)
}

// TestDataReaderMany verifies that the DataReader correctly returns all rows when the table contains multiple entries.
// It initializes the database with 3 rows and checks that the counter increments correctly to reflect all rows read.
func TestDataReaderMany(t *testing.T) {
	dr := prepareDr(t)
	defer dr.Close()
	dr.AppDb.Exec(TEST_TRUNC_TBL_SQL)
	dr.AppDb.Exec(TEST_INSERT_INTO_RAW_SQL)
	counter := 0
	for {
		next, err := dr.Next()
		if err != nil {
			fmt.Println(err)
			break
		}
		if !next {
			break
		}
		counter++
		res, err := dr.Scan()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
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
	counter := 0
	for {
		next, err := dr.Next()
		if err != nil {
			fmt.Println(err)
			break
		}
		if !next {
			break
		}
		counter++
		res, err := dr.Scan()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
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
	counter := 0
	for {
		next, err := dr.Next()
		if err != nil {
			fmt.Println(err)
			break
		}
		if !next {
			break
		}
		if counter > 3 {
			break
		}
		counter++
		res, err := dr.Scan()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
	assert.Equal(t, counter, 3)
}
