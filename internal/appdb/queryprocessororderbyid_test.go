// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestQueryProcessorOrderById verifies that the QueryProcessorOrderByID correctly appends the "WHERE id > {{id}}" and
// "ORDER BY id" clauses to the SQL query. It initializes the query processor, sets a limit, and asserts that the
// resulting query includes the expected clauses with the correct value for the {{id}} placeholder. The test checks
// that the value of the {{id}} placeholder is updated after calling SetValue with a new value.
func TestQueryProcessorOrderById(t *testing.T) {
	qp := QueryProcessorOrderByID{Query: "SELECT * FROM table WHERE id > {{id}} ORDER BY id LIMIT 10"}
	qp.InitQuery()
	actual := qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE id > 0 ORDER BY id LIMIT 10", actual)
	qp.SetValue("id", 10)
	actual = qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE id > 10 ORDER BY id LIMIT 10", actual)
}
