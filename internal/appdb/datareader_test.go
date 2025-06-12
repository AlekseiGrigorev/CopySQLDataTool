package appdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const DSN = "root:root@tcp(127.0.0.1:3306)/test"
const TBL_NAME = "test_table"
const TRUNC_TBL_SQL = "TRUNCATE TABLE " + TBL_NAME
const INSERT_INTO = "INSERT INTO "
const VALUES_123 = " VALUES (1), (2), (3)"
const SELECT_FROM = "SELECT * FROM "
const SELECT_TBL_SQL = SELECT_FROM + TBL_NAME + ";"

func prepareDb() DataReader {
	dr := DataReader{
		AppDb: &AppDb{Driver: "mysql", Dsn: DSN},
		Limit: 1,
		Query: SELECT_TBL_SQL,
	}
	dr.Open()
	return dr
}

func prepareDbOrderById() DataReader {
	dr := DataReader{
		AppDb:         &AppDb{Driver: "mysql", Dsn: DSN},
		Limit:         1,
		Query:         SELECT_FROM + TBL_NAME + " WHERE id > {{id}} ORDER BY id LIMIT 1;",
		Type:          TYPE_ORDERBYID,
		ExecutionTime: 3,
	}
	dr.AppDb.Open()
	return dr
}

func TestPrepareQuery(t *testing.T) {
	query := SELECT_TBL_SQL
	dr := DataReader{
		Query: SELECT_TBL_SQL,
		Limit: 1,
		Type:  TYPE_SIMPLE,
	}
	result := dr.prepareQuery()
	assert.Equal(t, query, result)
}

func TestPrepareQueryLimitOffset(t *testing.T) {
	query := SELECT_FROM + TBL_NAME + " LIMIT 1 OFFSET 0;"
	dr := DataReader{
		Query: SELECT_TBL_SQL + " ",
		Limit: 1,
		Type:  TYPE_LIMIT_OFFSET,
	}
	result := dr.prepareQuery()
	assert.Equal(t, query, result)
}

func TestPrepareQueryOrderById(t *testing.T) {
	query := SELECT_FROM + TBL_NAME + " WHERE id > 0 ORDER BY id LIMIT 1;"
	dr := DataReader{
		Query: SELECT_FROM + TBL_NAME + " WHERE id > {{id}} ORDER BY id LIMIT 1;",
		Type:  TYPE_ORDERBYID,
	}
	result := dr.prepareQuery()
	assert.Equal(t, query, result)
}

func TestDataReaderEmpty(t *testing.T) {
	dr := prepareDb()
	defer dr.Close()
	dr.AppDb.Exec(TRUNC_TBL_SQL)
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

func TestDataReaderOne(t *testing.T) {
	dr := prepareDb()
	defer dr.Close()
	dr.AppDb.Exec(TRUNC_TBL_SQL)
	dr.AppDb.Exec(INSERT_INTO + TBL_NAME + " VALUES (1)")
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

func TestDataReaderMany(t *testing.T) {
	dr := prepareDb()
	defer dr.Close()
	dr.AppDb.Exec(TRUNC_TBL_SQL)
	dr.AppDb.Exec(INSERT_INTO + TBL_NAME + VALUES_123)
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

func TestDataReaderManyOrderById(t *testing.T) {
	dr := prepareDbOrderById()
	defer dr.Close()
	dr.AppDb.Exec(TRUNC_TBL_SQL)
	dr.AppDb.Exec(INSERT_INTO + TBL_NAME + VALUES_123)
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

func TestDataReaderManyOrderByIdWrongType(t *testing.T) {
	dr := prepareDbOrderById()
	defer dr.Close()
	dr.AppDb.Exec(TRUNC_TBL_SQL)
	dr.AppDb.Exec(INSERT_INTO + TBL_NAME + VALUES_123)
	dr.Query = SELECT_FROM + TBL_NAME
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
