// Description: This package provides management features for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package app

import (
	"fmt"
	"os"
)

type FileProcessor struct {
	File *os.File
}

func (fp *FileProcessor) Write(buffer []string, data []any) error {
	for _, stmt := range buffer {
		_, err := fp.File.WriteString(stmt + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (fp *FileProcessor) GetProcessedMsg() string {
	return fmt.Sprint("Rows processed to file", fp.File.Name())
}
