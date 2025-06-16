// Description: This package provides file path management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appfilepath

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// AppFilePath represents a file path with additional metadata and utility methods for file path manipulation.
type AppFilePath struct {
	// Path to the file.
	Path string
}

// Returns the file path with the current date and time inserted
// into the file name. The date and time are formatted as "YYYYMMDD_HHMMSS".
// For example, if the original path is "/path/to/file.txt", the returned path
// would be "/path/to/file_20221231_000000.txt".
func (fp *AppFilePath) GetWithDateTime() string {
	dir, fileName := filepath.Split(fp.Path)
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)
	currentTime := time.Now().Format("20060102_150405")
	newFileName := fmt.Sprintf("%s_%s%s", name, currentTime, ext)
	return filepath.Join(dir, newFileName)
}
