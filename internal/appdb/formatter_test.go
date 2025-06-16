// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormatValueString tests the FormatValue method with a string value that contains special characters (backslash and single quote).
func TestFormatValueString(t *testing.T) {
	formatter := Formatter{}
	expected := "'test\\\\test''test'"
	actual := formatter.FormatValue("test\\test'test")
	assert.Equal(t, expected, actual)
}

// TestFormatRowValues tests the FormatRowValues method with a slice of values containing a string, an integer, and a float.
// It verifies that the method correctly formats the values into a string that can be used in an SQL statement.
func TestFormatRowValues(t *testing.T) {
	formatter := Formatter{}
	values := []any{"test\\test'test", 123, 123.456}
	fmt.Println(values)
	expected := "'test\\\\test''test', 123, 123.456"
	fmt.Println(expected)
	actual := formatter.FormatRowValues(values)
	fmt.Println(actual)
	assert.Equal(t, expected, actual)
}
