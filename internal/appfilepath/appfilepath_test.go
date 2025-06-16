// Description: This package provides file path management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appfilepath

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests the GetWithDateTime method of the AppFilePath type without a path.
// It tests that the method returns a file name with the current date and time inserted
// into the file name. The date and time are formatted as "YYYYMMDD_HHMMSS".
func TestGetWithDateTimeFileName(t *testing.T) {
	fp := AppFilePath{Path: "test.txt"}
	actual := fp.GetWithDateTime()
	fmt.Println(actual)
	assert.Regexp(t, `^test_\d{8}_\d{6}\.txt$`, actual)
}

// Tests the GetWithDateTime method of the AppFilePath type with a path.
// It tests that the method returns a file path with the current date and time inserted
// into the file name. The date and time are formatted as "YYYYMMDD_HHMMSS".
func TestGetWithDateTimePath(t *testing.T) {
	fp := AppFilePath{Path: "/var/log/test.txt"}
	actual := fp.GetWithDateTime()
	fmt.Println(actual)
	assert.Regexp(t, `^\Svar\Slog\Stest_\d{8}_\d{6}\.txt$`, actual)
}
