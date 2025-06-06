package appdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TBL_NAME = "test_table"
const TRUNC_TBL_SQL = "TRUNCATE TABLE " + TBL_NAME
const INSERT_INTO = "INSERT INTO "
const VALUES_123 = " VALUES (1), (2), (3)"
const SELECT_FROM = "SELECT * FROM "
const SELECT_TBL_SQL = SELECT_FROM + TBL_NAME + ";"

var DR DataReader

func prepareDb() {
	DR.Close()
	DR.AppDb = AppDb{Driver: "mysql", Dsn: "root:root@tcp(127.0.0.1:3306)/test"}
	DR.Limit = 1
	DR.Query = SELECT_TBL_SQL
	DR.AppDb.Open()
}

func prepareDbOrderById() {
	DR.Close()
	DR.AppDb = AppDb{Driver: "mysql", Dsn: "root:root@tcp(127.0.0.1:3306)/test"}
	DR.Limit = 1
	DR.Query = SELECT_FROM + TBL_NAME + " WHERE id > {{id}} ORDER BY id LIMIT 1;"
	DR.Type = "orderbyid"
	DR.ExecutionTime = 3
	DR.AppDb.Open()
}

func TestPrepareQuery(t *testing.T) {
	query := SELECT_FROM + TBL_NAME + " LIMIT 1 OFFSET 0;"
	m := DataReader{
		Query: SELECT_TBL_SQL + " ",
		Limit: 1,
	}
	result := m.prepareQuery()
	assert.Equal(t, query, result)
}

func TestPrepareQueryOrderById(t *testing.T) {
	query := SELECT_FROM + TBL_NAME + " WHERE id > 0 ORDER BY id LIMIT 1;"
	m := DataReader{
		Query: SELECT_FROM + TBL_NAME + " WHERE id > {{id}} ORDER BY id LIMIT 1;",
		Type:  "orderbyid",
	}
	result := m.prepareQuery()
	assert.Equal(t, query, result)
}

func TestDataReaderEmpty(t *testing.T) {
	prepareDb()
	DR.AppDb.Exec(TRUNC_TBL_SQL)
	counter := 0
	for {
		next, err := DR.Next()
		if err != nil {
			fmt.Println(err)
			break
		}
		if !next {
			break
		}
		counter++
		res, err := DR.Scan()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
	assert.Equal(t, counter, 0)
}

func TestDataReaderOne(t *testing.T) {
	prepareDb()
	DR.AppDb.Exec(TRUNC_TBL_SQL)
	DR.AppDb.Exec(INSERT_INTO + TBL_NAME + " VALUES (1)")
	counter := 0
	for {
		next, err := DR.Next()
		if err != nil {
			fmt.Println(err)
			break
		}
		if !next {
			break
		}
		counter++
		res, err := DR.Scan()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
	assert.Equal(t, counter, 1)
}

func TestDataReaderMany(t *testing.T) {
	prepareDb()
	DR.AppDb.Exec(TRUNC_TBL_SQL)
	DR.AppDb.Exec(INSERT_INTO + TBL_NAME + VALUES_123)
	counter := 0
	for {
		next, err := DR.Next()
		if err != nil {
			fmt.Println(err)
			break
		}
		if !next {
			break
		}
		counter++
		res, err := DR.Scan()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
	assert.Equal(t, counter, 3)
}

func TestDataReaderManyOrderById(t *testing.T) {
	prepareDbOrderById()
	DR.AppDb.Exec(TRUNC_TBL_SQL)
	DR.AppDb.Exec(INSERT_INTO + TBL_NAME + VALUES_123)
	counter := 0
	for {
		next, err := DR.Next()
		if err != nil {
			fmt.Println(err)
			break
		}
		if !next {
			break
		}
		counter++
		res, err := DR.Scan()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
	assert.Equal(t, counter, 3)
}

func TestDataReaderManyOrderByIdWrongType(t *testing.T) {
	prepareDbOrderById()
	DR.AppDb.Exec(TRUNC_TBL_SQL)
	DR.AppDb.Exec(INSERT_INTO + TBL_NAME + VALUES_123)
	DR.Query = SELECT_FROM + TBL_NAME
	DR.ExecutionTime = 1
	counter := 0
	for {
		next, err := DR.Next()
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
		res, err := DR.Scan()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
	assert.Equal(t, counter, 3)
}
