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
	TEST_DB_TBL_NAME             = "db.table"
	TEST_SELECT_FROM_DB_TBL_NAME = "SELECT * FROM " + TEST_DB_TBL_NAME
)

// TestGetFromTableName verifies that the GetFromTableName method extracts the table name correctly
// when the SQL query does not contain a WHERE clause and semicolon at the end.
func TestGetFromTableName(t *testing.T) {
	helper := SqlHelper{
		Sql: TEST_SELECT_FROM_DB_TBL_NAME,
	}
	expected := TEST_DB_TBL_NAME
	actual := helper.GetFromTableName()
	assert.Equal(t, expected, actual)
}

// TestGetFromTableNameWithSemicolon verifies that the GetFromTableName method extracts the table name
// correctly from a SQL query that ends with a semicolon.
func TestGetFromTableNameWithSemicolon(t *testing.T) {
	helper := SqlHelper{
		Sql: TEST_SELECT_FROM_DB_TBL_NAME + ";",
	}
	expected := TEST_DB_TBL_NAME
	actual := helper.GetFromTableName()
	assert.Equal(t, expected, actual)
}

// TestGetFromTableNameSpace verifies that the GetFromTableName method extracts the table name correctly
// even if the SQL query contains a WHERE clause with a space after the table name.
func TestGetFromTableNameSpace(t *testing.T) {
	helper := SqlHelper{
		Sql: TEST_SELECT_FROM_DB_TBL_NAME + " WHERE id = 1;",
	}
	expected := TEST_DB_TBL_NAME
	actual := helper.GetFromTableName()
	assert.Equal(t, expected, actual)
}

// TestStringify verifies that the SetStringify method correctly normalizes
// the SQL query string by replacing all whitespace characters with a single space.
func TestStringify(t *testing.T) {
	helper := SqlHelper{
		Sql: " SELECT\t* FROM\r\nTABLE    WHERE\n\t  ID\r= 1   \v\f; ",
	}
	helper.SetStringify()
	expected := "SELECT * FROM TABLE WHERE ID = 1 ;"
	actual := helper.Sql
	assert.Equal(t, expected, actual)
}
