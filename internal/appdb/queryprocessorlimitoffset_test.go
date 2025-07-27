// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestQueryProcessorLimitOffset verifies that the ProcessQuery method of the
// QueryProcessorLimitOffset struct correctly appends LIMIT and OFFSET clauses
// to the SQL query. It initializes the query processor, sets a limit, and
// asserts that the resulting query includes the expected LIMIT and OFFSET values.
// The test checks that the OFFSET is incremented by the limit value after each
// call to ProcessQuery, facilitating paging through results.
func TestQueryProcessorLimitOffset(t *testing.T) {
	qp := QueryProcessorLimitOffset{Query: "SELECT * FROM table ORDER BY id"}
	qp.InitQuery()
	qp.Limit = 10
	actual := qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table ORDER BY id LIMIT 10 OFFSET 0;", actual)
	actual = qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table ORDER BY id LIMIT 10 OFFSET 10;", actual)
}
