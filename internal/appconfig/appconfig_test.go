package appconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	configJSON = `{
    "config": {
        "source": {
            "driver": "mysql",
            "dsn": "test:test@tcp(localhost:3306)/test"
        },
        "dest": {
            "driver": "mysql",
            "dsn": "test:test@tcp(localhost:3306)/test"
        },
        "default_dataset": {
            "insert_command": "INSERT IGNORE INTO",
            "rows": 10000,
            "copy_to": "file",
            "query_type": "simple",
            "sql_statement": "prepared",
            "execution_time": 0
        }
    },
    "datasets": [
    	{
            "query": "SELECT * FROM db.test WHERE id = 1;",
            "table": ""
        }
	]
	}`
)

func TestEmptyTable(t *testing.T) {
	config := Config{}
	err := config.LoadConfigFromString(configJSON)
	if err != nil {
		t.Error("Error loading config:", err)
	}
	expected := "db.test"
	assert.Equal(t, expected, config.Datasets[0].Table)
}
