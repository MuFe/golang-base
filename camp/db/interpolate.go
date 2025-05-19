package db

import (
	"fmt"
	"strings"
	"time"
)

func InterpolateSQL(query string, args ...interface{}) string {
	for _, arg := range args {
		query = strings.Replace(query, "?", formatArg(arg), 1)
	}
	return query
}

func formatArg(arg interface{}) string {
	switch v := arg.(type) {
	case nil:
		return "NULL"
	case string:
		return "'" + escapeString(v) + "'"
	case []byte:
		return "'" + escapeString(string(v)) + "'"
	case time.Time:
		return "'" + v.Format("2006-01-02 15:04:05") + "'"
	case bool:
		if v {
			return "1"
		}
		return "0"
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	default:
		return fmt.Sprintf("'%v'", escapeString(fmt.Sprint(v)))
	}
}

func escapeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
