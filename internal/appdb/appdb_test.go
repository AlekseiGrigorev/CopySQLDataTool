// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Constants for testing.
const (
	TEST_DSN                   = "root:root@tcp(127.0.0.1:3306)/test"
	TEST_TBL_NAME              = "test_table"
	TEST_TRUNC_TBL_SQL         = "TRUNCATE TABLE " + TEST_TBL_NAME + ";"
	TEST_INSERT_INTO_SQL       = "INSERT INTO " + TEST_TBL_NAME + " VALUES (?), (?), (?);"
	TEST_INSERT_INTO_RAW_SQL   = "INSERT INTO " + TEST_TBL_NAME + " VALUES (1), (2), (3);"
	TEST_SELECT_FROM_SQL       = "SELECT * FROM " + TEST_TBL_NAME + ";"
	TEST_SELECT_FROM_LIMIT_SQL = "SELECT * FROM " + TEST_TBL_NAME + " LIMIT 1;"
	TEST_SELECT_FROM_COUNT_SQL = "SELECT COUNT(*) FROM " + TEST_TBL_NAME + ";"
)

// prepareDbAd initializes an AppDb instance with MySQL driver and DSN_AD.
// It opens the database connection and returns the AppDb instance.
// If opening the connection fails, it logs an error on the testing instance.
func prepareDbAd(t *testing.T) AppDb {
	ad := AppDb{
		Driver: "mysql",
		Dsn:    TEST_DSN,
	}
	err := ad.Open()
	if err != nil {
		t.Error(err)
	}
	return ad
}

// prepareDataAd truncates the table specified by TBL_NAME_AD in the provided AppDb instance
// and inserts a predefined set of values into it. It logs an error on the testing instance
// if the insertion fails.
func prepareDataAd(t *testing.T, ad AppDb) {
	truncateDataAd(t, ad)
	_, err := ad.Exec(TEST_INSERT_INTO_RAW_SQL)
	if err != nil {
		t.Error(err)
	}
}

// truncateDataAd truncates the table given by TBL_NAME_AD in the database connection provided by AppDb ad.
// It logs an error on the testing instance if the truncation fails.
func truncateDataAd(t *testing.T, ad AppDb) {
	_, err := ad.Exec(TEST_TRUNC_TBL_SQL)
	if err != nil {
		t.Error(err)
	}
}

// TestOpen tests the Open method of the AppDb struct.
// It creates an AppDb instance and calls the Open method on it.
// It then checks that the IsOpen method returns true,
// indicating that the database connection is open.
func TestOpen(t *testing.T) {
	ad := prepareDbAd(t)
	assert.True(t, ad.IsOpen())
}

// TestClose tests the Close method of the AppDb struct.
// It creates an AppDb instance, calls the Open method on it,
// calls the Close method, and checks that the IsOpen method returns false,
// indicating that the database connection is closed.
func TestClose(t *testing.T) {
	ad := prepareDbAd(t)
	ad.Close()
	assert.False(t, ad.IsOpen())
}

// TestExec verifies that the Exec method of the AppDb struct is functioning correctly.
// It initializes an AppDb instance and prepares the database with predefined data.
// The test asserts that the AppDb instance is not nil, indicating that the setup was successful.
func TestExec(t *testing.T) {
	ad := prepareDbAd(t)
	prepareDataAd(t, ad)
	assert.NotNil(t, ad)
}

// TestExecMultiple verifies that the ExecMultiple method of the AppDb struct can execute multiple SQL statements.
// It initializes an AppDb instance and prepares the database with predefined data.
// The test constructs a SQL command string with multiple statements and calls ExecMultiple,
// asserting that no errors occur during execution.
func TestExecMultiple(t *testing.T) {
	ad := prepareDbAd(t)
	prepareDataAd(t, ad)
	sql := TEST_SELECT_FROM_SQL + " " + TEST_SELECT_FROM_SQL
	err := ad.ExecMultiple(sql)
	assert.Nil(t, err)
}

// TestPrepareExec tests the PrepareExec method of the AppDb struct.
// It initializes an AppDb instance, truncates the table, and prepares an SQL INSERT statement.
// The test executes the statement with parameters and verifies that no error occurs
// and that the result is not nil, indicating successful execution.
func TestPrepareExec(t *testing.T) {
	ad := prepareDbAd(t)
	truncateDataAd(t, ad)
	sql := TEST_INSERT_INTO_SQL
	res, err := ad.PrepareExec(sql, 1, 2, 3)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

// TestQuery tests the Query method of the AppDb struct.
// It initializes an AppDb instance, truncates the table, and inserts data.
// The test calls Query with a SELECT statement and verifies that no error occurs and
// that the result is not nil, indicating successful execution.
func TestQuery(t *testing.T) {
	ad := prepareDbAd(t)
	prepareDataAd(t, ad)
	res, err := ad.Query(TEST_SELECT_FROM_SQL)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

// TestQueryRow verifies that the QueryRow method of the AppDb struct returns a valid Row instance
// when a query is executed with a LIMIT of 1. It initializes an AppDb instance, truncates the table,
// inserts data, and calls QueryRow with a SELECT statement. The test verifies that no error occurs
// and that the result is not nil, indicating successful execution.
func TestQueryRow(t *testing.T) {
	ad := prepareDbAd(t)
	prepareDataAd(t, ad)
	res, err := ad.QueryRow(TEST_SELECT_FROM_LIMIT_SQL)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

// TestGetScalar tests the GetScalar method of the AppDb struct.
// It initializes an AppDb instance, truncates the table, and inserts data.
// The test calls GetScalar with a SELECT statement and verifies that no error occurs and
// that the result is not nil, indicating successful execution.
func TestGetScalar(t *testing.T) {
	ad := prepareDbAd(t)
	res, err := ad.GetScalar(TEST_SELECT_FROM_COUNT_SQL)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
