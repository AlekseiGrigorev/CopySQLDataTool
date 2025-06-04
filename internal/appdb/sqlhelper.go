package appdb

import (
	"regexp"
)

type SqlHelper struct {
	Sql string
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
