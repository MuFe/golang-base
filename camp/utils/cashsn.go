package utils

import (
	"strconv"
	"time"
)

// GenerateCashSN 生成流水号
func GenerateCashSN() string {
	return strconv.Itoa(int(time.Now().Unix())) + Get4Code()
}
