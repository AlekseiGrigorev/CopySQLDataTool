// Description: This package provides management features for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWriteFile tests the Write method of the FileProcessor type.
// It creates a file, writes two strings to it, and checks that the
// message returned by GetProcessedMsg contains the name of the file.
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
