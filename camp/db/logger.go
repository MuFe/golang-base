package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/go-sql-driver/mysql"
	"github.com/ngrok/sqlmw"
	"time"
)

// Logger 实现 sqlmw.Interceptor 接口，拦截 Exec、Query、Prepare 等
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

// 重写 PrepareContext，返回自定义 Stmt 实现，拦截 Stmt.Exec 和 Stmt.Query
func (l Logger) PrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (driver.Stmt, error) {
	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &loggedStmt{
		Stmt:  stmt,
		query: query,
		log:   l,
	}, nil
}

// loggedStmt 包装 driver.Stmt，拦截 Exec 和 Query
type loggedStmt struct {
	driver.Stmt
	query string
	log   Logger
}

func (s *loggedStmt) Exec(args []driver.Value) (driver.Result, error) {
	start := time.Now()
	res, err := s.Stmt.Exec(args)
	duration := time.Since(start)

	// 转成 NamedValue 用 InterpolateSQL
	sqlStr := InterpolateSQL(s.query, valuesToInterface(args)...)
	logSQL("[STMT-EXEC]", sqlStr, duration, err)

	return res, err
}

func (s *loggedStmt) Query(args []driver.Value) (driver.Rows, error) {
	start := time.Now()
	rows, err := s.Stmt.Query(args)
	duration := time.Since(start)

	sqlStr := InterpolateSQL(s.query, valuesToInterface(args)...)
	logSQL("[STMT-QUERY]", sqlStr, duration, err)

	return rows, err
}

// 辅助：[]driver.NamedValue 转 []interface{}
func namedValueToInterface(args []driver.NamedValue) []interface{} {
	values := make([]interface{}, len(args))
	for i, v := range args {
		values[i] = v.Value
	}
	return values
}

// 辅助：[]driver.Value 转 []interface{}
func valuesToInterface(vals []driver.Value) []interface{} {
	values := make([]interface{}, len(vals))
	for i, v := range vals {
		values[i] = v
	}
	return values
}

// 输出日志函数，可替换成你自己的日志库
func logSQL(tag, sql string, duration time.Duration, err error) {
	slow := ""
	if duration > SlowThreshold {
		slow = "[SLOW]"
	}
	status := "[OK]"
	if err != nil {
		status = "[ERR]"
	}
	// 这里用标准log或你自定义logger
	logger.Printf("%s %s %s %v %s\n", tag, status, slow, duration, sql)

}

// RegisterWithLogging 注册带日志的 mysql 驱动
func RegisterWithLogging() {
	sql.Register("mysql-with-logger", sqlmw.Driver(mysql.MySQLDriver{}, Logger{}))
}
