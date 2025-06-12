package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteFile(t *testing.T) {
	filename := "test.txt"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	buffer := []string{"str1", "str2"}
	p := FileProcessor{File: file}
	p.Write(buffer, nil)
	actual := p.GetProcessedMsg()
	fmt.Println(actual)
	assert.Contains(t, actual, filename)
}
