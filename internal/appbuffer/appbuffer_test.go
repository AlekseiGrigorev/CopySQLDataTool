package appbuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	expected := "str1 str2"
	buff := AppBuffer{}
	actual := buff.AppendStr("str1").AppendStr("str2").Join(" ")
	assert.Equal(t, expected, actual)
}
