package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/ngrok/sqlmw"
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
	log.Printf("%s %s %s %v %s\n", tag, status, slow, duration, sql)
}

// wrappedDriver 包装 mysql.Driver，拦截 Open 以包装连接
type wrappedDriver struct {
	driver.Driver
}

func (d *wrappedDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return nil, err
	}

	// 包装连接，确保实现 ConnPrepareContext 和 ConnBeginTx
	var cpCtx driver.ConnPrepareContext
	if v, ok := conn.(driver.ConnPrepareContext); ok {
		cpCtx = v
	}

	var cbTx driver.ConnBeginTx
	if v, ok := conn.(driver.ConnBeginTx); ok {
		cbTx = v
	}

	return &wrappedConn{
		Conn:               conn,
		ConnPrepareContext: cpCtx,
		ConnBeginTx:        cbTx,
	}, nil
}

// wrappedConn 包装连接，确保实现 PrepareContext 和 BeginTx 拦截
type wrappedConn struct {
	driver.Conn
	driver.ConnPrepareContext
	driver.ConnBeginTx
}

func (c *wrappedConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if c.ConnPrepareContext != nil {
		stmt, err := c.ConnPrepareContext.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		return &loggedStmt{
			Stmt:  stmt,
			query: query,
		}, nil
	}
	return nil, driver.ErrSkip
}

func (c *wrappedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if c.ConnBeginTx != nil {
		tx, err := c.ConnBeginTx.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		return &wrappedTx{Tx: tx}, nil
	}
	// fallback
	return c.Conn.Begin()
}

// wrappedTx 包装事务，拦截 Commit 和 Rollback
type wrappedTx struct {
	driver.Tx
}

func (t *wrappedTx) Commit() error {
	log.Println("[TX COMMIT]")
	return t.Tx.Commit()
}

func (t *wrappedTx) Rollback() error {
	log.Println("[TX ROLLBACK]")
	return t.Tx.Rollback()
}

// loggedStmt 包装语句，拦截 Exec 和 Query
type loggedStmt struct {
	driver.Stmt
	query string
}

func (s *loggedStmt) Exec(args []driver.Value) (driver.Result, error) {
	start := time.Now()
	res, err := s.Stmt.Exec(args)
	duration := time.Since(start)

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

func (s *loggedStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now()
	execCtx := s.Stmt.(driver.StmtExecContext)
	res, err := execCtx.ExecContext(ctx, args)
	duration := time.Since(start)

	sqlStr := InterpolateSQL(s.query, namedValueToInterface(args)...)
	logSQL("[STMT-EXEC-CONTEXT]", sqlStr, duration, err)

	return res, err
}

func (s *loggedStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	start := time.Now()
	queryCtx := s.Stmt.(driver.StmtQueryContext)
	rows, err := queryCtx.QueryContext(ctx, args)
	duration := time.Since(start)

	sqlStr := InterpolateSQL(s.query, namedValueToInterface(args)...)
	logSQL("[STMT-QUERY-CONTEXT]", sqlStr, duration, err)

	return rows, err
}

func valuesToInterface(vals []driver.Value) []interface{} {
	values := make([]interface{}, len(vals))
	for i, v := range vals {
		values[i] = v
	}
	return values
}

// RegisterWithLogging 注册带日志拦截的自定义驱动
func RegisterWithLogging() {
	sql.Register("mysql-with-logger", &wrappedDriver{
		Driver: &mysql.MySQLDriver{},
	})
}
