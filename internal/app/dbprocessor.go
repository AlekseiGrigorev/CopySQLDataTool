// Description: This package provides management features for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package app

import (
	"copysqldatatool/internal/appdb"
	"fmt"
	"strings"
)

type DbProcessor struct {
	AppDb *appdb.AppDb
	Table string
}

func (db *DbProcessor) Write(buffer []string, data []any) error {
	_, err := db.AppDb.Exec(strings.Join(buffer, ""), data...)
	if err != nil {
		return fmt.Errorf("error writing to database: %w", err)
	}
	return nil
}

func (db *DbProcessor) GetProcessedMsg() string {
	return fmt.Sprint("Rows processed to table", db.Table)
}
