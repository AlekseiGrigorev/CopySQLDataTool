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
		Log:           nil,
		Rows:          rows,
		InsertCommand: INSERT_INTO,
		Table:         TBL_NAME_2,
		SqlStatement:  appdb.STATEMENT_RAW,
	}

	return &p
}

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

func TestWriteFileRp1(t *testing.T) {
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

func TestWriteFileRp2(t *testing.T) {
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

func TestWriteDbRp(t *testing.T) {
	p := prepareProcessor(&DbProcessor{AppDb: prepareDbRp(), Table: TBL_NAME_2}, 2)
	p.SqlStatement = appdb.STATEMENT_RAW
	err := p.Process()
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, "error writing to db")
		return
	}
	assert.Empty(t, err)
}

func TestWriteDbRpPrepared(t *testing.T) {
	p := prepareProcessor(&DbProcessor{AppDb: prepareDbRp(), Table: TBL_NAME_2}, 2)
	p.SqlStatement = appdb.STATEMENT_PREPARED
	err := p.Process()
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, "error writing to db")
		return
	}
	assert.Empty(t, err)
}
