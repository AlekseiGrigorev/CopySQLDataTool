// Description: This package provides log management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package applog

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	FILE     = "test.log"
	TEST_MSG = "Test message"
)

// Removes the test file if it exists after a test is run.
func cleanup() {
	if _, err := os.Stat(FILE); err == nil {
		os.Remove(FILE)
	}
}

// Tests the Info method of the AppLog type.
// It tests that the formatted string contains the provided message.
func TestInfo(t *testing.T) {
	log := AppLog{}
	s := log.String(TEST_MSG)
	fmt.Println(s)
	assert.Contains(t, s, TEST_MSG)
}

// Tests write to file of the AppLog type.
// It tests that the method writes the formatted string to the file, and
// that the string contains the provided message.
func TestWriteToFile(t *testing.T) {
	t.Cleanup(cleanup)
	file, err := os.Create(FILE)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	log := AppLog{
		File: file,
	}
	log.Info(TEST_MSG)
	log.Warn(TEST_MSG)
	log.Error(TEST_MSG)
	log.Ok(TEST_MSG)
	file.Seek(0, 0)
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	s := string(buffer[:n])
	assert.Contains(t, s, TEST_MSG)
}
