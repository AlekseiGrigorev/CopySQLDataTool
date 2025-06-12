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

type AppFilePath struct {
	Path string
}

func (fp *AppFilePath) GetWithDateTime() string {
	dir, fileName := filepath.Split(fp.Path)
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)
	currentTime := time.Now().Format("20060102_150405")
	newFileName := fmt.Sprintf("%s_%s%s", name, currentTime, ext)
	return filepath.Join(dir, newFileName)
}
