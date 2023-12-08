package xlog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	timeLocation *time.Location
	writer       io.Writer
	mutex        sync.Mutex
	skip         = 2
)

const (
	debugLevel = "DBG"
	infoLevel  = "INF"
	warnLevel  = "WRN"
	errorLevel = "ERR"

	timeFormart = "2006-01-02 15:04:05"
)

func init() {
	timeLocation = time.Now().Location()
	writer = os.Stdout
}

func output(s string) {
	writer.Write([]byte(s))
}

type msg struct {
	Level   string
	Time    time.Time
	File    string
	Line    int
	Func    string
	Message string
	Output  string
}

func (t *msg) format() {
	fileDir := strings.Split(t.File, "/")
	fileName := fileDir[len(fileDir)-1]
	msgList := strings.Split(t.Message, "\n")
	for i := range msgList {
		t.Output = t.Output + fmt.Sprintf("%s [%s] [%s():%s:%d] %s",
			t.Time.In(timeLocation).Format(timeFormart),
			t.Level,
			t.Func,
			fileName,
			t.Line,
			msgList[i]) + "\n"
	}
}

// Debug 打印日志
func Debug(params ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	message(debugLevel, skip, params...)
}

// Info 打印日志
func Info(params ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	message(infoLevel, skip, params...)
}

// Warn 打印警告
func Warn(params ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	message(warnLevel, skip, params...)
}

// Error 打印错误
func Error(params ...interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	return errors.New(message(errorLevel, skip, params...))
}

// Debugf 格式化打印调试
func Debugf(format string, params ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	message(debugLevel, skip, fmt.Sprintf(format, params...))
}

// Infof 格式化打印日志
func Infof(format string, params ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	message(infoLevel, skip, fmt.Sprintf(format, params...))
}

// Warnf 格式化打印警告
func Warnf(format string, params ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	message(warnLevel, skip, fmt.Sprintf(format, params...))
}

// Errorf 格式化打印错误并返回错误
func Errorf(format string, params ...interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	return errors.New(message(errorLevel, skip, fmt.Sprintf(format, params...)))
}

// ErrorP 打印错误
func ErrorP(params ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	message(errorLevel, skip, params...)
}

func message(level string, skip int, params ...interface{}) string {
	var message string
	var messageList []string
	for _, p := range params {
		messageList = append(messageList, fmt.Sprintf("%+v", p))
	}
	message = strings.Join(messageList, " ")
	functionID, _, _, _ := runtime.Caller(skip)
	function := runtime.FuncForPC(functionID)
	file, line := function.FileLine(functionID)
	m := msg{
		Level:   level,
		Time:    time.Now(),
		File:    file,
		Line:    line,
		Func:    function.Name(),
		Message: message,
	}
	m.format()
	output(m.Output)
	return message
}
