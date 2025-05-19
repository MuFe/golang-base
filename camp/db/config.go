package db

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	SlowThreshold           = 200 * time.Millisecond
	output        io.Writer = os.Stdout
	logger                  = log.New(output, "[MySQL] ", log.LstdFlags|log.Lshortfile)
)

func SetSlowThreshold(threshold time.Duration) {
	SlowThreshold = threshold
}

func SetLogOutput(w io.Writer) {
	output = w
	logger.SetOutput(w)
}
