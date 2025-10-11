// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestQueryProcessorBetween tests the QueryProcessorBetween struct by verifying that it correctly processes SQL queries with BETWEEN clauses
// and incrementing start and end values based on the provided step.
func TestQueryProcessorBetween(t *testing.T) {
	qp := QueryProcessorBetween{Query: "SELECT * FROM table WHERE field BETWEEN '{{start}}' AND '{{end}}' ORDER BY id"}
	qp.InitQuery()
	qp.Start = "2022-01-01 00:00:00"
	qp.End = "2022-01-02 00:00:00"
	qp.Step = "0h10m"
	actual := qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE field BETWEEN '2022-01-01 00:00:00' AND '2022-01-01 00:10:00' ORDER BY id", actual)
	actual = qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE field BETWEEN '2022-01-01 00:10:00' AND '2022-01-01 00:20:00' ORDER BY id", actual)
}

func TestQueryProcessorBetweenEnd(t *testing.T) {
	qp := QueryProcessorBetween{Query: "SELECT * FROM table WHERE field BETWEEN '{{start}}' AND '{{end}}' ORDER BY id"}
	qp.InitQuery()
	qp.Start = "2022-01-01 00:00:00"
	qp.End = "2022-01-01 00:10:00"
	qp.Step = "0h10m"
	qp.ProcessQuery()
	actual := qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE field BETWEEN '2022-01-01 00:10:00' AND '2022-01-01 00:10:00' ORDER BY id", actual)
	actual = qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE field BETWEEN '2022-01-01 00:20:00' AND '2022-01-01 00:10:00' ORDER BY id", actual)
	actual = qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE field BETWEEN '2022-01-01 00:20:00' AND '2022-01-01 00:10:00' ORDER BY id", actual)
}

// TestQueryProcessorBetweenInt verifies that the QueryProcessorBetween struct correctly processes SQL queries with BETWEEN clauses
// for integer fields and increments start and end values based on the provided step.
func TestQueryProcessorBetweenInt(t *testing.T) {
	qp := QueryProcessorBetween{Query: "SELECT * FROM table WHERE field BETWEEN '{{start}}' AND '{{end}}' ORDER BY id"}
	qp.InitQuery()
	qp.Start = "1"
	qp.End = "10"
	qp.Step = "1"
	actual := qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE field BETWEEN '1' AND '2' ORDER BY id", actual)
	actual = qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE field BETWEEN '2' AND '3' ORDER BY id", actual)
}

// TestQueryProcessorBetweenEndInt verifies that the QueryProcessorBetween struct correctly processes SQL queries with BETWEEN clauses for integer fields and increments start and end values based on the provided step when the end value is reached. It initializes the query processor, sets a limit, and asserts that the resulting query includes the expected clauses with the correct start and end values. The test checks that the start value is incremented by the step duration after each call, facilitating paging through results. If the current end value is after the end value specified in the query processor, the current end value is set to the end value.
func TestQueryProcessorBetweenEndInt(t *testing.T) {
	qp := QueryProcessorBetween{Query: "SELECT * FROM table WHERE field BETWEEN '{{start}}' AND '{{end}}' ORDER BY id"}
	qp.InitQuery()
	qp.Start = "1"
	qp.End = "4"
	qp.Step = "2"
	qp.ProcessQuery()
	actual := qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE field BETWEEN '3' AND '4' ORDER BY id", actual)
	actual = qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE field BETWEEN '5' AND '4' ORDER BY id", actual)
	actual = qp.ProcessQuery()
	assert.Equal(t, "SELECT * FROM table WHERE field BETWEEN '5' AND '4' ORDER BY id", actual)
}
