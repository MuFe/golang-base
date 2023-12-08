package xlog

import "fmt"

func NewCronLogger() *CronLogger {
	return &CronLogger{}
}

type CronLogger struct{}

func (t *CronLogger) Printf(s string, params ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	message(infoLevel, skip, fmt.Sprintf(s, params...))
}
