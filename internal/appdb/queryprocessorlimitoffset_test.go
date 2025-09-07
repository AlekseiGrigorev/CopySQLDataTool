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

// TestQueryProcessorLimitOffsetInitialLimit verifies that the ProcessQuery method of the
// QueryProcessorLimitOffset struct correctly appends LIMIT and OFFSET clauses
// to the SQL query when the Offset is set during initialization. It initializes the
// query processor, sets a limit, sets the offset to a non-zero value, and asserts
// that the resulting query includes the expected LIMIT and OFFSET values.
func TestQueryProcessorLimitOffsetInitialLimit(t *testing.T) {
	qp := QueryProcessorLimitOffset{Query: "SELECT * FROM table ORDER BY id"}
	qp.InitQuery()
	qp.Limit = 10
	qp.Offset = 10
	actual := qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table ORDER BY id LIMIT 10 OFFSET 10;", actual)
	actual = qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table ORDER BY id LIMIT 10 OFFSET 20;", actual)
}

// TestQueryProcessorLimitOffsetMaxOffset verifies that the ProcessQuery method of the
// QueryProcessorLimitOffset struct limits the offset to the MaxOffset value and resets
// the offset to 0 if it exceeds the MaxOffset. It initializes the query processor,
// sets a limit and a MaxOffset, and asserts that the resulting query includes the
// expected LIMIT and OFFSET values.
func TestQueryProcessorLimitOffsetMaxOffset(t *testing.T) {
	qp := QueryProcessorLimitOffset{Query: "SELECT * FROM table ORDER BY id"}
	qp.InitQuery()
	qp.Limit = 10
	qp.MaxOffset = 10
	actual := qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table ORDER BY id LIMIT 10 OFFSET 0;", actual)
	actual = qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table ORDER BY id LIMIT 10 OFFSET 10;", actual)
	actual = qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table ORDER BY id LIMIT 0 OFFSET 0;", actual)
}
