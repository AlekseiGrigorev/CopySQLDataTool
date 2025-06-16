// Description: This package provides configuration management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
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

// TestEmptyTable verifies that if the table name is empty in the config file,
// but a query is provided, the table name is extracted from the query.
func TestEmptyTable(t *testing.T) {
	config := Config{}
	err := config.LoadConfigFromString(configJSON)
	if err != nil {
		t.Error("Error loading config:", err)
	}
	expected := "db.test"
	assert.Equal(t, expected, config.Datasets[0].Table)
}
