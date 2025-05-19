package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/ngrok/sqlmw"
)

type Logger struct {
	sqlmw.NullInterceptor
}

func (l Logger) ExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now()
	res, err := conn.ExecContext(ctx, query, args)
	duration := time.Since(start)

	sqlStr := InterpolateSQL(query, namedValueToInterface(args)...)
	logSQL("[EXEC]", sqlStr, duration, err)

	return res, err
}

func (l Logger) QueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	start := time.Now()
	rows, err := conn.QueryContext(ctx, query, args)
	duration := time.Since(start)

	sqlStr := InterpolateSQL(query, namedValueToInterface(args)...)
	logSQL("[QUERY]", sqlStr, duration, err)

	return rows, err
}

func namedValueToInterface(args []driver.NamedValue) []interface{} {
	values := make([]interface{}, len(args))
	for i, v := range args {
		values[i] = v.Value
	}
	return values
}

func logSQL(tag, sql string, duration time.Duration, err error) {
	slow := ""
	if duration > SlowThreshold {
		slow = "[SLOW]"
	}
	status := "[OK]"
	if err != nil {
		status = "[ERR]"
	}
	logger.Printf("%s %s %s %v %s\n", tag, status, slow, duration, sql)
}

// RegisterWithLogging 注册 mysql-with-logger 驱动
func RegisterWithLogging() {
	sql.Register("mysql-with-logger", sqlmw.Driver(mysql.MySQLDriver{}, Logger{}))
}
