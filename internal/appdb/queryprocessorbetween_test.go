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
