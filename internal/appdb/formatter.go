package appdb

import (
	"fmt"
	"strconv"
	"strings"
)

type Formatter struct {
}

func (f *Formatter) FormatValue(val any) string {
	switch v := val.(type) {
	case []byte:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(strings.ReplaceAll(string(v), "'", "''"), "\\", "\\\\"))
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(strings.ReplaceAll(string(v), "'", "''"), "\\", "\\\\"))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

func (f *Formatter) FormatRowValues(values []any) string {
	var formattedValues []string

	for _, val := range values {
		formattedValues = append(formattedValues, f.FormatValue(val))
	}

	return strings.Join(formattedValues, ", ")
}

func (f *Formatter) BuildInsertPlaceholders(columnCount int) string {
	return strings.Repeat("?, ", columnCount-1) + "?"
}
