// Description: This package provides database management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"regexp"
	"strings"
)

// SqlHelper represents a helper struct for managing SQL queries.
// It contains the SQL query string that can be manipulated or analyzed.
type SqlHelper struct {
	// Sql is the SQL query string
	Sql string
}

// Stringify returns the SQL query string with all whitespace characters replaced by a single space.
func (sqlHelper *SqlHelper) Stringify() string {
	fields := strings.Fields(sqlHelper.Sql)
	return strings.Join(fields, " ")
}

// SetStringify sets the Sql field of the SqlHelper to the stringified version of the SQL query.
func (sqlHelper *SqlHelper) SetStringify() *SqlHelper {
	sqlHelper.Sql = sqlHelper.Stringify()
	return sqlHelper
}

// GetFromTableName extracts the table name from the SQL query stored in the SqlHelper.
// It uses a regular expression to find the first occurrence of the FROM clause and captures
// the table name that follows. Returns the table name if found, otherwise returns an empty string.
func (sqlHelper *SqlHelper) GetFromTableName() string {
	// Regular expression to match the first FROM clause and capture the table name
	re := regexp.MustCompile(`(?i)\bFROM\s+([^\s,;]+)`)
	match := re.FindStringSubmatch(sqlHelper.Sql)

	if len(match) > 1 {
		return match[1]
	}
	return ""
}
