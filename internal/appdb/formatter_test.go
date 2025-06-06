package appdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatValueString(t *testing.T) {
	formatter := Formatter{}
	expected := "'test\\\\test''test'"
	actual := formatter.FormatValue("test\\test'test")
	assert.Equal(t, expected, actual)
}

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
