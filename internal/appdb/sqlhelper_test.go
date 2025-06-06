package appdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	DB_TBL_NAME = "db.table"
)

func TestGetFromTableName(t *testing.T) {
	helper := SqlHelper{
		Sql: "SELECT * FROM db.table",
	}
	expected := DB_TBL_NAME
	actual := helper.GetFromTableName()
	assert.Equal(t, expected, actual)
}

func TestGetFromTableNameWithSemicolon(t *testing.T) {
	helper := SqlHelper{
		Sql: "SELECT * FROM db.table;",
	}
	expected := DB_TBL_NAME
	actual := helper.GetFromTableName()
	assert.Equal(t, expected, actual)
}

func TestGetFromTableNameSpace(t *testing.T) {
	helper := SqlHelper{
		Sql: "SELECT * FROM db.table WHERE id = 1;",
	}
	expected := DB_TBL_NAME
	actual := helper.GetFromTableName()
	assert.Equal(t, expected, actual)
}
