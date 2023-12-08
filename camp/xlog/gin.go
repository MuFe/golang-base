package xlog

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// NewGinLogger 新建logger
func NewGinLogger() *ginLogger {
	return &ginLogger{}
}

type ginLogger struct{}

func (t *ginLogger) Write(p []byte) (n int, err error) {
	message(infoLevel, skip, string(p))
	return len(p), err
}

// GinLogFormatter 格式
func GinLogFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s | %d | %s | %s | %s | %s | %s",
		"GIN",
		param.StatusCode,
		param.Method,
		param.Path,
		param.Latency,
		param.ClientIP,
		param.ErrorMessage,
	)
}
