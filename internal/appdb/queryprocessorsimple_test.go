// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestQueryProcessorSimple verifies that the QueryProcessorSimple struct correctly
// implements the QueryProcessorInterface by asserting that the ProcessQuery method
// returns the same query string as the one that was set, and that the InitQuery
// method does nothing. The test checks that the query string is not modified by
// either call.
func TestQueryProcessorSimple(t *testing.T) {
	qp := QueryProcessorSimple{Query: "SELECT * FROM table ORDER BY id"}
	qp.InitQuery()
	actual := qp.ProcessQuery()
	assert.Equal(t, qp.Query, actual)
	actual = qp.ProcessQuery()
	assert.Equal(t, qp.Query, actual)
}
