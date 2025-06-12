package app

import (
	"copysqldatatool/internal/appdb"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const DSN = "root:root@tcp(127.0.0.1:3306)/test"
const TBL_NAME = "test_table"
const TRUNC_TBL_SQL = "TRUNCATE TABLE " + TBL_NAME
const INSERT_INTO = "INSERT INTO "
const VALUES_123 = " VALUES (1), (2), (3)"
const INSERT_3 = INSERT_INTO + TBL_NAME + VALUES_123
const SELECT_FROM = "SELECT * FROM "
const SELECT_TBL_SQL = SELECT_FROM + TBL_NAME + ";"

func prepareDb() *appdb.AppDb {
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
	return &db
}

func TestWriteDb(t *testing.T) {
	buffer := []string{INSERT_INTO + TBL_NAME + " VALUES (1)", ", (2)", ", (3)", ";"}
	p := DbProcessor{
		AppDb: prepareDb(),
		Table: TBL_NAME,
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

func TestWriteDbPreparedStatement(t *testing.T) {
	buffer := []string{INSERT_INTO + TBL_NAME + " VALUES (?)", ", (?)", ", (?)", ";"}
	p := DbProcessor{
		AppDb: prepareDb(),
		Table: TBL_NAME,
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
