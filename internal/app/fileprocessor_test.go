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

// TestGetProcessedMsg tests the GetProcessedMsg method of the FileProcessor type when the File field is not set.
// It calls the GetProcessedMsg method and checks that the returned message contains the string "file is not set".
func TestGetProcessedMsgNilFile(t *testing.T) {
	p := FileProcessor{}
	actual := p.GetProcessedMsg()
	expected := "file is not set"
	assert.Contains(t, actual, expected)
}

// TestWriteNilFile tests the Write method of the FileProcessor type when the File field is not set.
// It calls the Write method with a buffer and nil data, and checks that the returned error
// is "file is not set".
func TestWriteNilFile(t *testing.T) {
	buffer := []string{"str1", "str2"}
	p := FileProcessor{}
	err := p.Write(buffer, nil)
	assert.Equal(t, err, fmt.Errorf("file is not set"))
}

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
