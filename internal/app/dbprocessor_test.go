// Description: This package provides management features for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package app

import (
	"copysqldatatool/internal/appdb"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Constants for testing.
const (
	DSN            = "root:root@tcp(127.0.0.1:3306)/test"
	TBL_NAME       = "test_table"
	TRUNC_TBL_SQL  = "TRUNCATE TABLE " + TBL_NAME
	INSERT_INTO    = "INSERT INTO "
	VALUES_123     = " VALUES (1), (2), (3)"
	INSERT_3       = INSERT_INTO + TBL_NAME + VALUES_123
	SELECT_FROM    = "SELECT * FROM "
	SELECT_TBL_SQL = SELECT_FROM + TBL_NAME + ";"
)

// prepareDb returns a new AppDb instance with the database connection opened.
// It truncates the table given by TBL_NAME after opening the database connection.
// It returns nil in case of any error during the process.
func prepareDb(t *testing.T) *appdb.AppDb {
	db := appdb.AppDb{
		Driver: "mysql",
		Dsn:    DSN,
	}
	err := db.Open()
	if err != nil {
		t.Error(err)
		return nil
	}
	_, err = db.Exec(TRUNC_TBL_SQL)
	if err != nil {
		t.Error(err)
		return nil
	}
	return &db
}

// TestGetProcessedMsgNilDb tests the GetProcessedMsg method of the DbProcessor type when the AppDb field is not set.
// It calls the GetProcessedMsg method and checks that the returned message contains the string "db is not set".
func TestGetProcessedMsgNilDb(t *testing.T) {
	p := DbProcessor{}
	actual := p.GetProcessedMsg()
	expected := "db is not set"
	assert.Contains(t, actual, expected)
}

// TestWriteNilDb tests the Write method of the DbProcessor type when the AppDb field is not set.
// It calls the Write method with a buffer and nil data, and checks that the returned error
// is "db is not set".
func TestWriteNilDb(t *testing.T) {
	buffer := []string{"str1", "str2"}
	p := DbProcessor{}
	err := p.Write(buffer, nil)
	assert.Equal(t, err, fmt.Errorf("db is not set"))
}

// TestWriteDb tests the Write method of the DbProcessor struct with a raw SQL statement.
// It uses a buffer with 4 strings: "INSERT INTO test_table VALUES (1)", ", (2)", ", (3)", and ";".
// It then calls the Write method with the buffer and an empty slice of any.
// The test checks that the Write method does not return an error and that the GetProcessedMsg method
// returns a message that contains the table name.
func TestWriteDb(t *testing.T) {
	buffer := []string{INSERT_INTO + TBL_NAME + " VALUES (1)", ", (2)", ", (3)", ";"}
	p := DbProcessor{
		AppDb:     prepareDb(t),
		TableName: TBL_NAME,
	}
	err := p.Write(buffer, []any{})
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, "error writing to database")
		return
	}
	actual := p.GetProcessedMsg()
	fmt.Println(actual)
	assert.Contains(t, actual, TBL_NAME)
}

// TestWriteDbPreparedStatement tests the Write method of the DbProcessor struct with a prepared SQL statement.
// It uses a buffer with 4 strings: "INSERT INTO test_table VALUES (?)", ", (?)", ", (?)", and ";".
// It then calls the Write method with the buffer and a slice of any with 3 elements.
// The test checks that the Write method does not return an error and that the GetProcessedMsg method
// returns a message that contains the table name.
func TestWriteDbPreparedStatement(t *testing.T) {
	buffer := []string{INSERT_INTO + TBL_NAME + " VALUES (?)", ", (?)", ", (?)", ";"}
	p := DbProcessor{
		AppDb:     prepareDb(t),
		TableName: TBL_NAME,
	}
	err := p.Write(buffer, []any{1, 2, 3})
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, "error writing to database")
		return
	}
	actual := p.GetProcessedMsg()
	fmt.Println(actual)
	assert.Contains(t, actual, TBL_NAME)
}
