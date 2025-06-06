package applog

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TEST_MSG = "Test message"
)

func TestInfo(t *testing.T) {
	log := AppLog{}
	s := log.String(TEST_MSG)
	fmt.Println(s)
	assert.Contains(t, s, TEST_MSG)
}

func TestWriteToFile(t *testing.T) {
	file, err := os.Create("test.log")
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
