// Description: This package provides string buffer management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appbuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests the Join method to ensure it correctly concatenates strings with a separator.
func TestJoin(t *testing.T) {
	expected := "str1 str2"
	buff := AppBuffer{}
	actual := buff.AppendStr("str1").AppendStr("str2").Join(" ")
	assert.Equal(t, expected, actual)
}
