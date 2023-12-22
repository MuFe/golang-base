package util

import (
	"fmt"
	"strconv"
	"strings"
)

// mysql in写法
func MysqlStringInUtils(buf *strings.Builder, list []int64, preString string) {
	list = RemoveDupsInt64(list)
	if len(list) == 1 && list[0] == 0 {
		return
	}
	for k, info := range list {
		if len(list) == 1 {
			buf.WriteString(preString)
			buf.WriteString(" =")
			buf.WriteString(strconv.FormatInt(info, 10))
		} else {
			if k == 0 {
				buf.WriteString(preString)
				buf.WriteString(" in (")
			}
			buf.WriteString(strconv.FormatInt(info, 10))
			if k < len(list)-1 {
				buf.WriteString(",")
			} else {
				buf.WriteString(") ")
			}
		}

	}
}

// mysql in写法
func MysqlStringInUtilsWithZero(buf *strings.Builder, list []int64, preString string) {
	list = RemoveDupsInt64(list)
	for k, info := range list {
		if len(list) == 1 {
			buf.WriteString(preString)
			buf.WriteString(" =")
			buf.WriteString(strconv.FormatInt(info, 10))
		} else {
			if k == 0 {
				buf.WriteString(preString)
				buf.WriteString(" in (")
			}
			buf.WriteString(strconv.FormatInt(info, 10))
			if k < len(list)-1 {
				buf.WriteString(",")
			} else {
				buf.WriteString(") ")
			}
		}

	}
}

// mysql in写法
func MysqlInUtils(buf *strings.Builder, list []string, preString string) {
	list = RemoveDupsString(list)
	if len(list) == 1 && list[0] == "" {
		return
	}
	for k, info := range list {
		if k == 0 {
			buf.WriteString(preString)
		}
		buf.WriteString(fmt.Sprintf("'%s'", info))
		if k < len(list)-1 {
			buf.WriteString(",")
		} else {
			buf.WriteString(") ")
		}
	}
}
